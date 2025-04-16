package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type QueryAddr struct {
	Connection net.PacketConn
	Client     net.Addr
	Query      Q
}

type IncomingRequest struct {
	Connection net.PacketConn
	Client     net.Addr
	Buffer     []byte
	BuffSize   int
}

func sendOkay(client net.PacketConn, client_addr net.Addr, buffer []byte) QueryAddr {
	log.Println("Sending okay to ", client_addr)

	var q Q

	err := q.UnmarshalBinary(buffer)
	if err != nil {
		log.Println("Server Handler error decoding query: ", err)
	}

	log.Println("Deccoded Q: ", q)

	response := q.Response("okay")

	response_buffer, err := response.MarshallBinary()
	if err != nil {
		log.Println("Server Handler error encoding response: ", err)
	}
	client.WriteTo(response_buffer, client_addr)

	log.Println("Server Sent: ", response)

	return QueryAddr{Connection: client, Client: client_addr, Query: q}
}

func sendNotify(data QueryAddr) R {
	log.Println("Sending Notify to ", data.Client)

	response := data.Query.Response("Notify")

	log.Println("Responding with ", response)

	var response_buffer bytes.Buffer
	encoder := gob.NewEncoder(&response_buffer)

	encoder.Encode(response)
	data.Connection.WriteTo(response_buffer.Bytes(), data.Client)

	return response
}

func runTwoResponseServer() {
	log.Println("Starting two response server")
	requests := make(chan IncomingRequest)
	send_notify := make(chan QueryAddr)
	Okays := make(chan R)
	go func() {
		req := <-requests
		d := sendOkay(req.Connection, req.Client, req.Buffer[:req.BuffSize])
		log.Println("Adding ", d, " to Notifies")
		send_notify <- d
	}()
	go func() {
		n := <-send_notify
		log.Println("From Notifies channel: ", n)
		Okays <- sendNotify(n)
	}()
	go func() {
		r := <-Okays
		log.Println("From Okays channel: ", r)
	}()
	listener, err := net.ListenPacket("udp", "127.0.0.1:50000")
	if err != nil {
		log.Println("Server Error setting up listener: ", err)
	}

	log.Println("Listening on ", listener.LocalAddr())
	//defer listener.Close()

	for {
		buff := make([]byte, 1024)
		bytes_read, client_addr, err := listener.ReadFrom(buff)
		log.Println("Number of bytes received: ", bytes_read)
		if err != nil {
			log.Println("Server Error accepting request: ", err)
		}
		requests <- IncomingRequest{Connection: listener, Client: client_addr, Buffer: buff, BuffSize: bytes_read}
		log.Println("Received request from ", client_addr)
	}

}
