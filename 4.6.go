package main

import (
	"fmt"
)

func main() {
	s := []byte("  这是 \n   什么 \r  东西")
	s = removeSpace(s)
	fmt.Printf("%s\n", string(s))
}

func removeSpace(str []byte) []byte{
	flag := 0
	for _, b := range str {
		if isSpace(b) {
			if flag > 0 && isSpace(str[flag-1]) {
				continue
			}
			b = ' '
		}
		str[flag] = b
		flag++
	}
	return str[:flag]
}

func isSpace(chr byte) bool {
	return chr == ' ' || chr == '\t' || chr == '\n' || chr == '\r'
}