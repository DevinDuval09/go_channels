package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

func sendSubscribe(server net.UDPAddr, id int) net.PacketConn {
	//connect to server
	conn, err := net.DialUDP("udp", nil, &server)
	//create subscribe query
	var q = Q{Message: "Subscribe", Number: id}

	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)

	encoder.Encode(q)

	if err != nil {
		fmt.Println("Client Error while connecting to server: ", err)
	}

	conn.Write(b.Bytes())
	return conn
}

func acceptNotify(conn net.PacketConn) {
	var r R
	rbuf := make([]byte, 1024)
	network := net.Buffers{rbuf}

	_, _, err := conn.ReadFrom(rbuf)
	if err != nil {
		fmt.Println("Client error reading Notify: ", err)
	}

	decoder := gob.NewDecoder(&network)

	decoder.Decode(&r)

	fmt.Println("Received Notify: ", r)
}
