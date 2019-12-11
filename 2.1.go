package main

import (
	"./tempconv"
	"fmt"
)

func main() {
	var c tempconv.Celsius = 100
	fmt.Printf("Celsius = %s, kevin = %s\n", c, tempconv.CToK(c))		
}