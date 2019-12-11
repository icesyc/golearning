package main

import (
	"fmt"
	"os"
	"crypto/sha256"
)

func main() {
	for _, arg := range os.Args[1:] {
		fmt.Printf("%x\n", sha256.Sum256([]byte(arg)))
	}
}