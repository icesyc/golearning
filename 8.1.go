package main

import (
	"net"
	"io"
	"log"
	"time"
	"flag"
)

var port = flag.String("port", "8000", "listen port")
var offset = flag.Int("offset", 2, "timezone offset")

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", "localhost:" + *port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept();
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("client incoming %s\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Add(time.Duration(*offset) * time.Hour).Format("15:04:05\n"))
		if err != nil {
			log.Print(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}