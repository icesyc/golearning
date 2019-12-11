package main

import (
	"net"
	"fmt"
	"bufio"
)

type Client struct{
	Name string
	Ch chan string
}
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
	client := Client{"", make(chan string)}
	scanner := bufio.NewScanner(conn)
	fmt.Fprintln(conn, "please input your name:")
	scanner.Scan()
	client.Name = scanner.Text()

	messages <- client.Name + " has arrived"
	entering <- client
	go clientWriter(conn, client)
	for scanner.Scan() {
		msg := scanner.Text()	
		messages <- client.Name + ": " + msg
	}

	leaving <- client
	messages <- client.Name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, client Client) {
	for msg := range client.Ch {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster() {
	var clients = make(map[Client]bool)
	for {
		select{
		case msg := <-messages:
			for client, _ := range clients {
				client.Ch <- msg
			}
		case client := <-entering:
			clients[client] = true
		case client := <-leaving:
			delete(clients, client)
			close(client.Ch)
		}
	}
}