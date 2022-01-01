// Package vsohashcuda implements the paged VSO-Hash using CUDA.
// See github.com/peterebden/vso-hash for more details and a CPU-based implementation.
package vsohashcuda

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"runtime"
	"sync"
	"unsafe"
)

/*
#include "sha256.h"
#cgo LDFLAGS: -L. -lsha256
*/
import "C"

// Size is the number of bytes of the output hash.
const Size = sha256.Size + 1

// BlockSize is the size of the VSO-Hash blocks. It's defined as always being 2MB.
const BlockSize = 2 * 1024 * 1024

// PageSize is the size of the pages within each block. They're always 64kb
const PageSize = 64 * 1024

const pagesPerBlock = BlockSize / PageSize

const seed = "VSO Content Identifier Seed"

// initOnce guards the sha256_preinit function
var initOnce sync.Once

// A Hasher holds GPU-related resources necessary for computing SHA-256 hashes.
type Hasher struct {
	jobs    **C.SHA256_job
	numJobs int
	maxSize int
}

// New creates a new Hasher.
// It will calculate up to the given number of hashes in parallel.
//
// It is safe for parallel use. Typically one would create a single Hasher and share it in order to make best use of resources.
//
// It does not implement the builtin hash.Hash type since we have many cases where errors can occur, whereas hash.Hash is
// not allowed to; it's also hard for us to parallelise as we'd like against that interface.
func New(parallelism int) (*Hasher, error) {
	initOnce.Do(func() {
		C.sha256_preinit()
	})

	jobs := C.sha256_alloc_jobs(C.int(parallelism), C.int(PageSize))
	if jobs == nil {
		return nil, cudaError(C.sha256_last_error())
	}
	h := &Hasher{
		jobs:    jobs,
		numJobs: parallelism,
	}
	runtime.SetFinalizer(h, finalize)
	return h, nil
}

// finalize is run when the GC will collect a Hasher; it frees GPU-side resources.
func finalize(h *Hasher) {
	C.sha256_free_jobs(h.jobs, C.int(h.numJobs))
}

// SHA256Sums calculates SHA256 hashes in parallel for the given byte slices.
// They can at most be 64kb since that is our page size.
func (h *Hasher) SHA256Sums(data [][]byte) ([][sha256.Size]byte, error) {
	if len(data) > h.numJobs {
		return nil, fmt.Errorf("Provided byte slices exceeds max number of jobs (max %d, got %d)", h.numJobs, len(data))
	}
	for i, datum := range data {
		if len(datum) > PageSize {
			return nil, fmt.Errorf("Provided byte slice exceeds max page size of 64kb (was %d)", len(datum))
		} else if C.sha256_init_job(h.jobs, C.int(i), C.SHA256_data(unsafe.Pointer(&datum[0])), C.int(len(datum))) != 0 {
			return nil, cudaError(C.sha256_last_error())
		}
	}
	C.sha256_run(h.jobs, C.int(len(data)))
	if err := C.sha256_last_error(); err != nil {
		return nil, cudaError(err)
	}
	ret := make([][sha256.Size]byte, len(data))
	for i := range ret {
		C.sha256_copy_digest(h.jobs, C.int(i), C.SHA256_data(unsafe.Pointer(&ret[i])))
	}
	return ret, nil
}

// cudaError returns the last error that happened on the CUDA side.
func cudaError(err *C.char) error {
	return fmt.Errorf("CUDA error: %s", C.GoString(err))
}

// Hash calculates the paged VSO hash over the given reader.
func (h *Hasher) Hash(r io.Reader) ([]byte, error) {
	pageHashes := make(chan []byte, h.numJobs)
	blockHashes := make(chan []byte, 10)
	vsoHash := make(chan []byte)
	go func() {
		// Aggregate page digests together into block digests
		var buf bytes.Buffer
		const max = pagesPerBlock * sha256.Size
		buf.Grow(max)
		for dg := range pageHashes {
			buf.Write(dg)
			if buf.Len() == max {
				sum := sha256.Sum256(buf.Bytes())
				blockHashes <- sum[:]
				buf.Reset()
			}
		}
		// Remember to hash the last block, if there is overhang
		if buf.Len() > 0 {
			sum := sha256.Sum256(buf.Bytes())
			blockHashes <- sum[:]
		}
		close(blockHashes)
	}()

	// Aggregate block digests together at the top level
	go func() {
		var buf bytes.Buffer
		buf.Grow(2*Size + 1)
		buf.WriteString(seed)
		buf.Write(<-blockHashes)
		for hash := range blockHashes {
			buf.WriteByte(0) // not the last block
			sum := sha256.Sum256(buf.Bytes())
			buf.Reset()
			buf.Write(sum[:])
			buf.Write(hash)
		}
		buf.WriteByte(1) // this is the last block
		sum := sha256.Sum256(buf.Bytes())
		vsoHash <- append(sum[:], 0)
	}()

	page := make([]byte, PageSize)
	eof := false
	first := true
	for !eof {
		i := 0
		for ; i < h.numJobs; i++ {
			if n, err := r.Read(page); err != nil {
				if err == io.EOF {
					eof = true
					break
				}
				return nil, err
			} else if result := C.sha256_init_job(h.jobs, C.int(i), C.SHA256_data(unsafe.Pointer(&page[0])), C.int(n)); result != 0 {
				return nil, cudaError(C.sha256_last_error())
			}
		}
		if i == 0 {
			if first {
				// The empty file is a special case; it has one page hash of zero bytes.
				sum := sha256.Sum256(nil)
				pageHashes <- sum[:]
			}
			continue
		}
		first = false
		// We've either filled up all our jobs or hit the end of the file - either way we need to run a new set of hashes.
		C.sha256_run(h.jobs, C.int(i))
		if err := C.sha256_last_error(); err != nil {
			return nil, cudaError(err)
		}
		for j := 0; j < i; j++ {
			digest := make([]byte, sha256.Size)
			C.sha256_copy_digest(h.jobs, C.int(j), C.SHA256_data(unsafe.Pointer(&digest[0])))
			pageHashes <- digest
		}
	}
	close(pageHashes)
	return <-vsoHash, nil
}
