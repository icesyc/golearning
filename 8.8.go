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

func scan(c net.Conn, line chan string) {
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		line <- scanner.Text()
	}
	if scanner.Err() != nil {
		log.Print(scanner.Err())
	}
}
func handleConnection(c net.Conn) {
	fmt.Printf("%s incomding\n", c.RemoteAddr())
	defer c.Close()
	line := make(chan string)
	go scan(c, line)
	timeout := 5 * time.Second
	abort := time.NewTimer(timeout)
	for {
		select {
		case <-abort.C:
			return
		case text := <-line:
			abort.Reset(timeout)
			go echo(c, text, 1 * time.Second)
		}
	}
}