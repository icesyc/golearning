package main

import (
	"os"
	"fmt"
)

func main() {
	fmt.Printf("joined: %s\n", Join(",", os.Args[1:]...))
}

func Join(sep string, arr ...string) string{
	var res string
	for i, s := range arr {
		if i > 0 {
			res += sep + s
		} else {
			res += s
		}
	}
	return res
}