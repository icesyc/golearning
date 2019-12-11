package main

import (
	"net"
	"os"
	"strings"
	"fmt"
	"io"
	"bufio"
)


func main() {
	var cities []string
	var Connections = map[string]net.Conn{}
	for _, pair := range os.Args[1:] {
		pairArray := strings.Split(pair, "=")
		cities = append(cities, pairArray[0])
		Connections[pairArray[0]], _ = net.Dial("tcp", pairArray[1])
	}

	for _, city := range cities {
		fmt.Printf("%-20s", city)
	}
	fmt.Printf("\n")
	for {
		fmt.Printf("\r")
		for _, city := range cities {
			conn := Connections[city]	
			reader := conn.(io.Reader)
			breader := bufio.NewReader(reader)
			time, err := breader.ReadString('\n')
			if err != nil {
				time = fmt.Sprintf("%s", err)
			}
			fmt.Printf("%-20s", strings.TrimSpace(time))
		}
	}
}

