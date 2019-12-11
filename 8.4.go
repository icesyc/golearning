package main

import (
	"fmt"
	"net"
	"log"
	"bufio"
	"time"
	"strings"
	"sync"
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
	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		text := scanner.Text()
		go func(text string){
			defer wg.Done()
			echo(c, text, 1 * time.Second)
		}(text)
	}
	wg.Wait()	
	if scanner.Err() != nil {
		log.Print(scanner.Err())
	}
	time.Sleep(2*time.Second)
	tcpConn := c.(*net.TCPConn)
	tcpConn.CloseWrite()
}