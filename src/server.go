package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

func handle(client net.PacketConn, client_addr net.Addr, buffer []byte) {
	var q Q
	network := net.Buffers{buffer}
	decoder := gob.NewDecoder(&network)
	err := decoder.Decode(&q)
	if err != nil {
		fmt.Println("Server Handler error decoding query: ", err)
	}

	response := q.Response("server response")

	var response_buffer bytes.Buffer
	encoder := gob.NewEncoder(&response_buffer)

	encoder.Encode(response)
	client.WriteTo(response_buffer.Bytes(), client_addr)

}

func runTwoResponseServer() {
	listener, err := net.ListenPacket("udp", ":50000")
	if err != nil {
		fmt.Println("Server Error setting up listener: ", err)
	}

	defer listener.Close()

	for {
		buff := make([]byte, 1024)
		bytes_read, client_addr, err := listener.ReadFrom(buff)
		if err != nil {
			fmt.Println("Server Error accepting request: ", err)
		}
		go handle(listener, client_addr, buff[:bytes_read])
	}

}
