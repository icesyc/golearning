package main

import (
	"net"
	"os"
	"log"
	"io"
)

func main() {
	c, err := net.Dial("tcp", ":8000")	
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	go io.Copy(os.Stdout, c)
	io.Copy(c, os.Stdin)
}