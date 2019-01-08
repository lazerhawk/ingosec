package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

	if !listen && len(target) > 0 && port > 0 {
		clientSender()
	}

	// we are going to listen and possibly upload, execute commands, or drop a shell
	if listen {
		serverLoop()
	}
}

func serverLoop() {

	ln, err := net.Listen("tcp", target+":"+strconv.Itoa(port))
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
	// check for upload
	if len(uploadDestination) > 0 {
		buffer := ""
		data := make([]byte, 1024)
		for {
			size, err := conn.Read(data)
			if err != nil {
				break
			} else {
				buffer += string(data[:size])
			}
		}
		err := ioutil.WriteFile(uploadDestination, []byte(buffer), 0777)
		if err != nil {
			fmt.Fprintf(conn, "Failed to save file to %s\r\n", uploadDestination)
		}
		fmt.Fprintf(conn, "Successfully saved file to %s\r\n", uploadDestination)
	}
	// check for command execution
	if len(execute) > 0 {
		cmdOutput := runCommand(execute)
		fmt.Fprintf(conn, cmdOutput)
	}
	// loop if cmd shell was requested
	if command {
		for {
			fmt.Fprintf(conn, "<GNC:#> ")
			data := make([]byte, 1024)
			buffer := ""
		ReadLoop:
			for {
				size, err := conn.Read(data)
				if err != nil {
					panic(err)
				}
				buffer += string(data[:size])
				for _, v := range data {
					if v == '\n' {
						break ReadLoop
					}
				}
			}
			cmdOutput := runCommand(buffer)
			fmt.Fprintf(conn, cmdOutput)
		}
	}
}

func runCommand(command string) string {
	command = strings.Trim(command, " ")
	command = strings.Trim(command, "\n")

	cmd := exec.Command(command)
	cmdOut, err := cmd.Output()
	if err != nil {
		cmdOut = []byte("Failed to execute command: " + command)
	}
	return string(cmdOut)
}

func readRaw() string {
	// read in the buffer
	scanner := bufio.NewScanner(os.Stdin)

	// send out the buffer to the clientSender
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func clientSender() {
	buffer := ""
	conn, err := net.Dial("tcp", target+":"+strconv.Itoa(port))
	defer conn.Close()
	if err != nil {
		fmt.Printf("Error accepting a connection, %v", err)
	}
	for {
		buffer = readRaw()
		if len(buffer) > 0 {
			fmt.Fprintln(conn, buffer)
		}
		recvLen := 1
		response := ""
		buf := make([]byte, 4096)
		for {
			if recvLen > 0 {
				bufLen, _ := bufio.NewReader(conn).Read(buf)
				recvLen = bufLen
				response += string(buf[:bufLen])
				if recvLen < 4096 {
					break
				}
			} else {
				break
			}
		}
		fmt.Printf(response)
	}
}
