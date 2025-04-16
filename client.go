package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

var sub_count = 0

func sendSubscribe(server net.UDPAddr, messages chan string, out chan net.PacketConn) {
	log.Println("Client sending subscribe...")
	//create subscribe query
	message := <-messages
	sub_count++
	q := Q{Number: sub_count, Message: message}
	log.Println("Client created query ", q)
	//connect to server
	conn, err := net.DialUDP("udp", nil, &server)
	if err != nil {
		log.Println("Client Error while connecting to server: ", err)
	}

	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)

	err = encoder.Encode(q)
	if err != nil {
		log.Println("Client error while encoding: ", err)
	}

	log.Println("Encoded q")

	conn.Write(b.Bytes())
	log.Println("Sent query")
	out <- conn
}

func handleOk(conns chan net.PacketConn, out chan net.PacketConn) {
	log.Println("Client handling ok...")
	conn := <-conns
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
	log.Println("Client handling notify...", count)
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

	fmt.Println("Received: ", r.Message)

	q := Q{Number: count, Message: "ok"}

	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)

	encoder.Encode(q)

	log.Println("Received Notify: ", r)
	conn.Close()
}
