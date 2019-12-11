package main

import (
	"os"
	"fmt"
	"bufio"
)

func main(){
	counts := make(map[string]int)
	countFiles := make(map[string]string)

	args := os.Args[1:]
	if len(args) == 0 {
		countLines(os.Stdin, counts, countFiles, "stdin")
	} else {
		for _, file := range args {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				continue
			}
			countLines(f, counts, countFiles, file)
		}
	}

	for line, n := range counts {
		fmt.Printf("%d\t%s\t%s\n", n, line, countFiles[line])
	}
}

func countLines(f *os.File, counts map[string]int, countFiles map[string]string, file string){
	input := bufio.NewScanner(f)
	for input.Scan() {
		text := input.Text()
		counts[text]++
		countFiles[text] += file + " "
	}	
}