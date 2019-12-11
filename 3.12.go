package main

import (
	"fmt"
	"os"
)
func main() {
	if len(os.Args) != 3 {
		fmt.Printf("wrong parameters num, must be 2\n")
		os.Exit(1)
	}
	res := "not same"
	if sameChars(os.Args[1], os.Args[2]) {
		res = "same"
	}
	fmt.Printf("[%s] compare [%s] is %s\n", os.Args[1], os.Args[2], res)
}

func sameChars(s1, s2 string) bool {
	c1 := countChars(s1)
	c2 := countChars(s2)
	for chr, cnt := range c1 {
		if c2[chr] != cnt {
			return false;
		}
	}
	for chr, cnt := range c2 {
		if c1[chr] != cnt {
			return false;
		}
	}
	return true;
}

func countChars(s string) map[rune]int {
	m := make(map[rune]int)
	for _, c := range s {
		m[c]++
	}
	return m
}