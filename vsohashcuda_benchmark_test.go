package vsohashcuda_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/peterebden/vso-hash-cuda"
)

func BenchmarkVSOHashCUDA(b *testing.B) {
	// Need quite a bit of data to give our hashes a chance to shine.
	const size = 2 * 1024 * 1024 * 1024

	data := make([]byte, size)
	for i := uint64(0); i < size; i += 8 {
		binary.LittleEndian.PutUint64(data[i:], i)
	}
	b.ResetTimer()

	for _, parallelism := range []int{100, 200, 400, 800, 1600, 2400, 4800} {
		b.Run(fmt.Sprintf("Parallel%d", parallelism), func(b *testing.B) {
			h, err := vsohashcuda.New(parallelism)
			if err != nil {
				b.Fatalf("Failed to create hasher: %s", err)
			}
			b.ResetTimer()
			start := time.Now()
			for i := 0; i < b.N; i++ {
				if _, err := h.Hash(bytes.NewReader(data)); err != nil {
					b.Fatalf("Failed to hash data: %s", err)
				}
			}
			b.ReportMetric(float64(size*b.N)/(1024*1024*time.Since(start).Seconds()), "MB/s")
		})
	}
}
