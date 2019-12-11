package main

import (
	"fmt"
	"sort"
)

func main() {
	str := String("abccba")
	fmt.Printf("%v\n", isPalindrome(str))
}

type String string
func (s String) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s String) Len() int {
	return len(s)
}
func (s String) Swap(i, j int) {
	return	
}
func equal(i, j int, s sort.Interface) bool {
	return !s.Less(i, j) && !s.Less(j, i)
}
func isPalindrome(s sort.Interface) bool {
	slen := s.Len()
	for i := 0; i < slen / 2; i++ {
		j := slen - i - 1
		if !equal(i, j, s) {
			return false
		}	
	}
	return true
}
