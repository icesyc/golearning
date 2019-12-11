package main

import (
	"fmt"
)

type s struct {
	name string
}
func main() {
	fmt.Printf("%v\n", returnInt())
}

func returnInt() (r interface{}) {
	defer func() {
		if p := recover(); p != nil {
			r = p
		}
	}()
	panic("non zero result returned")
}