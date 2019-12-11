package main

import (
	"fmt"
)

func main() {
	s := []string{"apple", "apple", "orange", "banana", "banana", "ice", "ice"}
	s = removeDup(s)
	fmt.Printf("%v\n", s)
}

func removeDup(str []string) []string{
	idx := 0
	for _, s := range str {
		if s == str[idx] {
			continue
		}
		idx++
		str[idx] = s
	}
	return str[:idx+1]
}