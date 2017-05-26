package main

import (
	// "bufio"
	"fmt"
	"net"
	"os"
	// "strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	ln, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer ln.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	for {
		reqLen, err := conn.Read(buf)
		if reqLen == 0 {
			return
		}
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		fmt.Printf("%3d>> %s\n", reqLen, buf[0:reqLen])
	}
}

//
// func main() {
//
// 	fmt.Println("Launching server...")
//
// 	// listen on all interfaces
// 	ln, _ := net.Listen("tcp", ":8081")
//
// 	// accept connection on port
// 	conn, _ := ln.Accept()
//
// 	// run loop forever (or until ctrl-c)
// 	for {
// 		// will listen for message to process ending in newline (\n)
// 		message, err := bufio.NewReader(conn).ReadString('\n')
// 		if err != nil {
// 			fmt.Printf("%v", err)
// 		}
// 		if len(message) > 0 {
// 			// output message received
// 			fmt.Print("Message Received:", string(message))
// 			// sample process for string received
// 			newmessage := strings.ToUpper(message)
// 			// send new string back to client
// 			conn.Write([]byte(newmessage + "\n"))
//
// 		}
// 	}
// }
