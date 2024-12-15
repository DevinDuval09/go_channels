package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type QueryAddr struct {
	Connection net.PacketConn
	Client     net.Addr
	Query      Q
}

func sendOkay(client net.PacketConn, client_addr net.Addr, buffer []byte) QueryAddr {
	var q Q
	network := net.Buffers{buffer}
	decoder := gob.NewDecoder(&network)
	err := decoder.Decode(&q)
	if err != nil {
		fmt.Println("Server Handler error decoding query: ", err)
	}

	response := q.Response("okay")

	var response_buffer bytes.Buffer
	encoder := gob.NewEncoder(&response_buffer)

	encoder.Encode(response)
	client.WriteTo(response_buffer.Bytes(), client_addr)

	return QueryAddr{Connection: client, Client: client_addr, Query: q}
}

func sendNotify(data QueryAddr) R {
	response := data.Query.Response("Notify")

	var response_buffer bytes.Buffer
	encoder := gob.NewEncoder(&response_buffer)

	encoder.Encode(response)
	data.Connection.WriteTo(response_buffer.Bytes(), data.Client)

	return response
}

func runTwoResponseServer() {
	listener, err := net.ListenPacket("udp", ":50000")
	if err != nil {
		fmt.Println("Server Error setting up listener: ", err)
	}
	Notifies := make(chan QueryAddr)
	Okays := make(chan R)

	defer listener.Close()

	for {
		buff := make([]byte, 1024)
		bytes_read, client_addr, err := listener.ReadFrom(buff)
		if err != nil {
			fmt.Println("Server Error accepting request: ", err)
		}
		go func() { Notifies <- sendOkay(listener, client_addr, buff[:bytes_read]) }()
		go func() {
			Okays <- sendNotify(<-Notifies)
		}()
		go func() {
			r := <-Okays
			fmt.Println(r)
		}()
	}

}
