package main

import (
	"net"
	"fmt"
	"log"
	"io"
	"os"
	"strings"
	"bufio"
	"io/ioutil"
	"time"
	"math/rand"
	"strconv"
	"bytes"
)

var rootDir, _ = os.Getwd()
func main() {
	server, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
	server.Close()
}

func handleDataConn(port int, dataChan chan string) {
	host := ":" + strconv.Itoa(port)
	server, err := net.Listen("tcp", host)
	if err != nil {
		log.Print(err)
		return
	}
	defer server.Close()
	conn, err := server.Accept()
	if err != nil {
		log.Print(err)
		return
	}
	command := <-dataChan
	switch command {
	case "retr", "list": 
		out := <-dataChan
		io.WriteString(conn, out + "\n")
	case "stor":
		data, _ := ioutil.ReadAll(conn)
		fmt.Printf("%d bytes\n", len(data))
		dataChan<- string(data)
	}
	conn.Close()
}

func handleConn(c net.Conn) {
	defer c.Close()
	fmt.Fprintf(c, "220 (ice ftp server 1.0)\n")
	reader := c.(io.Reader)
	scanner := bufio.NewScanner(reader)
	var dataChan = make(chan string)
	var mode = ""
	for {
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Print(err)
			return
		}
		fmt.Printf("scanned %v\n", scanner.Text())
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		cmd := strings.ToLower(fields[0])
		var args []string
		if len(fields) > 1 {
			args = fields[1:]
		}
		switch cmd {
		case "close", "exit", "quit":
			fmt.Fprintf(c, "Bye\n")
			return
		case "user": {
			fmt.Fprintf(c, "331 Please specify the password\n")
			continue
		}
		case "pass": {
			fmt.Fprintf(c, "230 Login successful.\n")
			continue
		}
		case "syst": {
			fmt.Fprintf(c, "215 UNIX Type: L8\n")
			continue
		}
		case "feat": {
			fmt.Fprintf(c, "211-Features:\n PASV\n EPSV\n211 End\n")
			continue
		}
		case "pwd": {
			fmt.Fprintf(c, "257 %q\n", rootDir)
			continue
		}
		case "epsv": {
			s := rand.NewSource(time.Now().UnixNano())
			rnd := rand.New(s)
			port := rnd.Intn(64511) + 1024
			fmt.Fprintf(c, "229 Entering Extended Passive Mode (|||%d|).\n", port)
			go handleDataConn(port, dataChan)
			mode = "pasv"
		}
		case "pasv":
			addr := c.LocalAddr().String()
			tmp := strings.Split(addr, ":")
			ip := strings.Join(strings.Split(tmp[0], "."), ",")
			s := rand.NewSource(time.Now().UnixNano())
			rnd := rand.New(s)
			port := rnd.Intn(64511) + 1024
			p1 := port / 256
			p2 := port % 256
			fmt.Fprintf(c, "227 Entering Passive Mode (%s,%d,%d).\n", ip, p1, p2)
			go handleDataConn(port, dataChan)
			mode = "pasv"
		case "type": 
			//only use binary mode
			fmt.Fprintf(c, "200 Switching to Binary mode.\n")
		case "list":
			if mode != "pasv" {
				fmt.Fprintf(c, "425 Use PORT or PASV first.\n")
				continue
			}
			files, err := ioutil.ReadDir(".")
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			var buf bytes.Buffer
			for _, file := range files {
				fmt.Fprintf(&buf, "%s\t", file.Name())
			}
			fmt.Fprintf(&buf, "\n")
			dataChan<- cmd
			dataChan<- buf.String()
			fmt.Fprintf(c, "150 Here comes the directory listing.\n")
			fmt.Fprintf(c, "226 Directory send OK.\n")
			mode = ""
		case "cwd":
			if len(args) == 0 {
				os.Chdir(rootDir)
				continue
			}
			os.Chdir(args[1])
			fmt.Printf("250 Directory successfully changed.\n")
		case "retr":
			if mode != "pasv" {
				fmt.Fprintf(c, "425 Use PORT or PASV first.\n")
				continue
			}
			if len(args) == 0 {
				fmt.Fprintf(c, "550 Failed to open file.\n")
				continue
			}
			f, err := os.Stat(args[0])
			if err != nil {
				fmt.Fprintf(c, "550 Failed to open file.\n")
				continue
			}
			fmt.Fprintf(c, "150 Opening BINARY mode data connection for %s (%d bytes).\n", f.Name(), f.Size())
			content, err := ioutil.ReadFile(args[0])
			dataChan<- cmd
			dataChan<- string(content)
			fmt.Fprintf(c, "226 Transfer complete.\n")                                    
		case "stor":
			fmt.Fprintf(c, "150 Ok to send data.\n")
			dataChan<- cmd
			fileName := args[0]
			content := <-dataChan
			ioutil.WriteFile(fileName, []byte(content), 0644)
			fmt.Fprintf(c, "226 Transfer complete.\n")                                    
		default:
			fmt.Fprintf(c, "550 error command.\n")
		}
	}
}