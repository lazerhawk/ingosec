package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")

	if err != nil {

	} else {
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		status, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(status)

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			fmt.Fprintf(conn, scanner.Text()+"\n")
			status, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Printf("From server: %s", status)
		}

	}
}
