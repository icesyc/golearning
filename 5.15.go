package main

import (
	"os"
	"fmt"
	"strconv"
)

func main() {
	var ints []int
	for _, arg := range os.Args[1:] {
		i, _ := strconv.Atoi(arg)
		ints = append(ints, i)
	}
	maxInt := max(ints...)
	minInt := min(ints...)
	fmt.Printf("max is %d\n", maxInt)
	fmt.Printf("min is %d\n", minInt)
}

func max(args ...int) int {
	if len(args) == 0 {
		return 0
	}
	m := args[0]
	for _, arg := range args {
		if arg > m {
			m = arg
		}
	}
	return m
}

func min(args ...int) int {
	if len(args) == 0 {
		return 0
	}
	m := args[0]
	for _, arg := range args {
		if arg < m {
			m = arg
		}
	}
	return m
}