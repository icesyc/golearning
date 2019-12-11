package main

import (
	"fmt"
	"strings"
)

func main() {
	s := expand("this is foo", quote)
	fmt.Println(s)
}
func quote(s string) string {
	return "(" + s + ")"
}

func expand(s string, f func(string) string) string {
	if f != nil {
		return strings.Replace(s, "foo", f("foo"), -1)
	}
	return s
}