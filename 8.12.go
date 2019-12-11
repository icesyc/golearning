package main

import (
	"net"
	"fmt"
	"bufio"
)

type Client chan string
var messages = make(chan string)
var entering = make(chan Client)
var leaving = make(chan Client)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	client := make(Client)
	who := conn.RemoteAddr().String()
	messages <- who + " has arrived"
	entering <- client
	go clientWriter(conn, client)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()	
		messages <- who + ": " + msg
	}

	leaving <- client
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch Client) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster() {
	var clients = make(map[Client]bool)
	for {
		select{
		case msg := <-messages:
			for client, _ := range clients {
				client <- msg
			}
		case client := <-entering:
			clients[client] = true
		case client := <-leaving:
			delete(clients, client)
			close(client)
		}
	}
}