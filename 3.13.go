package main

import "fmt"

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
)
func main() {
	fmt.Printf("%v %v %v\n", KB, MB, GB)
}