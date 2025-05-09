package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	fmt.Println("Executing main")
	server_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:50000")
	if err != nil {
		log.Println("Driver error resolving server address: ", err)
	}
	log.Println("Setting up client...")
	out_subscribes := make(chan string)
	incoming_ok := make(chan net.PacketConn)
	incoming_notify := make(chan net.PacketConn)
	fmt.Println("Set up channels")
	go sendSubscribe(*server_addr, out_subscribes, incoming_ok)
	go handleOk(incoming_ok, incoming_notify)
	go handleNotify(incoming_notify)
	log.Println("Client made channels...")
	time.Sleep(2 * time.Second)
	for i := 0; i < 4; i++ {
		log.Println("Sending subscribe number ", i+1)
		out_subscribes <- "my message"
	}
}
