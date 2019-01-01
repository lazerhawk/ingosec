package main

import (
	"flag"
	"fmt"
)

var execute = ""
var target = "localhost"
var uploadDestination = ""
var upload = false
var listen = false
var command = false
var port = 8080

func main() {
	targetPtr := flag.String("t", target, "target_host")
	portPtr := flag.Int("p", port, "target_port")
	listenPtr := flag.Bool("l", listen, "listen on [host]:[port] for incoming connections")
	commandPtr := flag.Bool("c", command, "initialize a command shell")
	uploadPtr := flag.String("u", uploadDestination, "upon receiving a connection upload and write a file to [destination]")
	executePtr := flag.String("e", execute, "execute the given file upon receiving a connection")

	flag.Parse()

	execute = *executePtr
	target = *targetPtr
	uploadDestination = *uploadPtr
	if uploadDestination != "" {
		upload = true
	}
	listen = *listenPtr
	command = *commandPtr
	port = *portPtr

	fmt.Printf("Target Host: %s\n", target)
	fmt.Printf("Target Port: %v\n", port)
	fmt.Printf("Listening Mode: %v\n", listen)
}
