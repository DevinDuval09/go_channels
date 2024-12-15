package main

import (
	"fmt"
	"net"
)

func main() {
	server_addr, err := net.ResolveUDPAddr("udp", ":50000")
	if err != nil {
		fmt.Println("Driver error resolving server address: ", err)
	}

	Okays := make(chan net.PacketConn)

	for i := 0; i < 3; i++ {
		go func() { Okays <- sendSubscribe(*server_addr, i) }()
		go acceptNotify(<-Okays)
	}
}
