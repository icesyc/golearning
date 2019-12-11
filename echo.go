package main

import (
	"fmt"
	"net"
	"log"
	"bufio"
	"time"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func echo(c net.Conn, str string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(str))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", str)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(str))
}

func handleConnection(c net.Conn) {
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		go echo(c, scanner.Text(), 1 * time.Second)
	}
	if scanner.Err() != nil {
		log.Print(scanner.Err())
	}
	c.Close()
}