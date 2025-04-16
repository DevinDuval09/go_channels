package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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
	Data       []byte
}

func sendOkay(client net.PacketConn, client_addr net.Addr, q Q) QueryAddr {
	fmt.Println("Server sending okay to ", client_addr)

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
	fmt.Println("Server sending Notify to ", data.Client)

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
		log.Println("Server starting request handler...")
		req := <-requests
		log.Println("Recieved request")
		reader := bytes.NewReader(req.Data)
		var q Q
		err := binary.Read(reader, binary.BigEndian, q)
		if err != nil {
			log.Println("Got error parsing buffer: ", err)
		}
		d := sendOkay(req.Connection, req.Client, q)
		log.Println("Adding ", d, " to Notifies")
		send_notify <- d
	}()
	go func() {
		log.Println("Server starting notify sender...")
		n := <-send_notify
		log.Println("From Notifies channel: ", n)
		Okays <- sendNotify(n)
	}()
	go func() {
		log.Println("Server starting Okay sender...")
		r := <-Okays
		log.Println("From Okays channel: ", r)
	}()
	listener, err := net.ListenPacket("udp", "127.0.0.1:50000")
	if err != nil {
		log.Println("Server Error setting up listener: ", err)
	}

	log.Println("Listening on ", listener.LocalAddr())
	defer listener.Close()

	for {
		buff := make([]byte, 1024)
		bytes_read, client_addr, err := listener.ReadFrom(buff)
		data := buff[:bytes_read]
		log.Println("received data: ", data)
		if err != nil {
			log.Println("Server Error accepting request: ", err)
		}
		requests <- IncomingRequest{Connection: listener, Client: client_addr, Data: data}
		log.Println("Received request from ", client_addr)
	}

}
