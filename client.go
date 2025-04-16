package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

func sendSubscribe(server net.UDPAddr, ids chan int, messages chan string, out chan net.PacketConn) {
	log.Println("Sending subscribe")
	//connect to server
	conn, err := net.DialUDP("udp", nil, &server)
	if err != nil {
		log.Println("Client Error while connecting to server: ", err)
	}
	//create subscribe query
	log.Println("sendSubscribe connected to address")
	id := <-ids
	message := <-messages
	q := Q{Number: id, Message: message}

	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)

	encoder.Encode(q)

	//log.Println("Encoded q")

	log.Println("Sending subscribe")
	conn.Write(b.Bytes())
	log.Println("Client sent: ", q)
	out <- conn
}

func handleOk(conns chan net.PacketConn, out chan net.PacketConn) {
	log.Println("Waiting on connection...")
	conn := <-conns
	log.Println("Received connection")
	rbuf := make([]byte, 1024)
	_, _, err := conn.ReadFrom(rbuf)
	if err != nil {
		log.Println("Client error reading Okay: ", err)
	}

	out <- conn
}

var count int = 0

func handleNotify(conns chan net.PacketConn) {
	conn := <-conns
	log.Println("Handling notify ", count)
	count++
	var r R
	rbuf := make([]byte, 1024)
	network := net.Buffers{rbuf}
	_, _, err := conn.ReadFrom(rbuf)
	if err != nil {
		log.Println("Client error reading Notify: ", err)
	}

	decoder := gob.NewDecoder(&network)

	decoder.Decode(&r)

	q := Q{Number: count, Message: "ok"}

	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)

	encoder.Encode(q)

	log.Println("Received Notify: ", r)
	conn.Close()
}
