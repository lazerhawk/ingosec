package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error starting server, %v", err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting a connection, %v", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(conn).ReadString('\n') //read until newline
		if err != nil {
			fmt.Printf("Error reading from connection, %v", err)
			return
		}
		fmt.Printf("Received from client: %s", data)
		temp := strings.TrimSpace(string(data))
		if temp == "STOP" {
			break
		}
		conn.Write([]byte(string(data)))
	}
	conn.Close()
}
