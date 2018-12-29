package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "google.com:80")

	if err != nil {

	} else {
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		status, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(status)
	}
}
