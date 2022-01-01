// Package main implements a minimal main designed to test the CUDA kernel.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/peterebden/vso-hash"
	"github.com/peterebden/vso-hash-cuda"
)

var inputs = [][]byte{
	[]byte("The most merciful thing in the world, I think, is the inability of the human mind to correlate all its contents."),
	[]byte("We live on a placid island of ignorance in the midst of black seas of infinity, and it was not meant that we should voyage far."),
	[]byte("The sciences, each straining in its own direction, have hitherto harmed us little; but some day the piecing together of dissociated knowledge will open up such terrifying vistas of reality, and of our frightful position therein, that we shall either go mad from the revelation or flee from the deadly light into the peace and safety of a new dark age."),
	[]byte("Theosophists have guessed at the awesome grandeur of the cosmic cycle wherein our world and human race form transient incidents."),
}

func main() {
	fmt.Printf("CPU SHA256:\n")
	for i, input := range inputs {
		sum := sha256.Sum256(input)
		fmt.Printf("Sentence %d: %s\n", i, hex.EncodeToString(sum[:]))
	}

	fmt.Printf("GPU SHA256:\n")
	h, err := vsohashcuda.New(10) // Length is arbitrary but longer than we need.
	if err != nil {
		log.Fatalf("%s", err)
	}
	sums, err := h.SHA256Sums(inputs)
	if err != nil {
		log.Fatalf("%s", err)
	}
	for i, sum := range sums {
		fmt.Printf("Sentence %d: %s\n", i, hex.EncodeToString(sum[:]))
	}

	fmt.Printf("CPU VSO-Hash:\n")
	for i, input := range inputs {
		sum := vsohash.Sum(input)
		fmt.Printf("Sentence %d: %s\n", i, hex.EncodeToString(sum[:]))
	}

	fmt.Printf("GPU VSO-Hash:\n")
	for i, input := range inputs {
		sum, err := h.Hash(bytes.NewReader(input))
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Printf("Sentence %d: %s\n", i, hex.EncodeToString(sum[:]))
	}
}
