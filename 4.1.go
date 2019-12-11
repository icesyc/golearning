package main

import (
	"fmt"
	"crypto/sha256"
)

func main() {
	s1, s2 := "x", "X"
	sha1 := sha256.Sum256([]byte(s1))
	sha2 := sha256.Sum256([]byte(s2))
	fmt.Printf("%x vs %x diff bit count: %d\n", sha1, sha2, bitDiffCount(sha1, sha2))
}

func bitDiffCount(sha1, sha2 [32]uint8) int {
	n := 0
	for i := range sha1 {
		n += bitNum(sha1[i] ^ sha2[i])
	}
	return n
}

func bitNum(b uint8) int{
	n := 0
	for b > 0 {
		b &= b - 1
		n++
	}
	return n
}