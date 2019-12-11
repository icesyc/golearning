package main

import (
	"os"
	"fmt"
	"unicode"
	"unicode/utf8"
	"bufio"
	"io"
)

func main() {
	letterCounts := make(map[rune]int)
	digitCounts := make(map[rune]int)
	utflen := make([]int, utf8.UTFMax + 1)
	invalid := 0

	in := bufio.NewReader(os.Stdin)

	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		if unicode.IsLetter(r) {
			letterCounts[r]++
		}else if unicode.IsDigit(r) {
			digitCounts[r]++
		}
		utflen[n]++
	}

	fmt.Printf("letter\tcounts\n")
	for r, n := range letterCounts {
		fmt.Printf("%q\t%d\n", r, n)
	}

	fmt.Printf("digit\tcounts\n")
	for r, n := range digitCounts {
		fmt.Printf("%q\t%d\n", r, n)
	}

	fmt.Printf("\nlength\tcounts\n")
	for l, n := range utflen {
		if n > 0 {
			fmt.Printf("%d\t%d\n", l, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}
