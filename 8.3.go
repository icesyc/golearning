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
	done := make(chan struct{})
	go func(){
		io.Copy(os.Stdout, c)
		log.Println("done")
		done<- struct{}{}
	}()
	go func(){
		io.Copy(c, os.Stdin)
	}()
	tcpConn := c.(*net.TCPConn)
	tcpConn.CloseWrite()
	<-done
}