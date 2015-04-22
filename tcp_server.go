package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// RfDataServer starts the server on specified host / port and return
// received RF data on parser channel. The host can be an IP or the host
// to listen on. parser channel is where the received data is sent.
// This is a blocking call, it starts an inifinite loop waiting on
// connections. Each connection is passed onto a goroutine to be handled
// asynchronously.
func RfDataServer(host string, port uint32, parser chan<- string) {
	laddr := fmt.Sprintf("%s:%d", host, port)

	sock, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Printf("Error connecting to %s: %s\n", laddr, err)
	}
	defer sock.Close()

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Printf("Error establishing connection: %s\n", err)
		}
		defer conn.Close()

		go receiveData(conn, parser)
	}
}

func receiveData(connection net.Conn, parser chan<- string) {
	var message string
	var err error

	for {
		message, err = bufio.NewReader(connection).ReadString('\n')
		if err != nil && err != io.EOF {
			log.Printf("Error reading from connection: %s\n", err)
			continue
		}

		if 0 == len(message) {
			continue
		}

		parser <- message
	}
}
