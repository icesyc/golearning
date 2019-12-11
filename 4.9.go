package main

import (
	"os"
	"fmt"
	"bufio"
)

func main() {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		w := scanner.Text()
		counts[w]++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("word\tcounts\n")
	for w, n := range counts {
		fmt.Printf("%s\t%d\n", w, n)
	}
}
