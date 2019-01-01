package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:80")

	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	} else {
		fmt.Fprintf(conn, "Hello UDP Server")
		_, err := bufio.NewReader(conn).Read(p)
		if err != nil {
			fmt.Printf("Some error reading, %v", err)
			return
		}
		fmt.Printf("%s\n", p)
	}
	conn.Close()
}
