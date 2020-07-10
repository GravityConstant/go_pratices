package main

import (
	"fmt"
	"net"
)

func main() {
	ipv6()
}

func ipv6() {
	uaddr, err := net.ResolveUDPAddr("udp", "[::1]:8080")
	if err != nil {
		fmt.Printf("get udp ipv6 error: %v\n", err)
	}
	conn, err := net.DialUDP("udp", nil, uaddr)
	if err != nil {
		fmt.Printf("connect ipv6 error: %v\n", err)
	}

	_, err = conn.Write([]byte("Hello server ipv6"))
	if err != nil {
		fmt.Printf("send data failed ipv6, err: %v\n", err)
		return
	}

	result := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(result)
	if err != nil {
		fmt.Printf("receive data failed, err: %v\n", err)
		return
	}
	fmt.Printf("receive from addr: %v, data: %v\n", remoteAddr, string(result[:n]))
}

func ipv4() {
	//connect server
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})

	if err != nil {
		fmt.Printf("connect failed, err: %v\n", err)
		return
	}

	//send data
	_, err = conn.Write([]byte("hello server!"))
	if err != nil {
		fmt.Printf("send data failed, err : %v\n", err)
		return
	}

	//receive data from server
	result := make([]byte, 4096)
	n, remoteAddr, err := conn.ReadFromUDP(result)
	if err != nil {
		fmt.Printf("receive data failed, err: %v\n", err)
		return
	}
	fmt.Printf("receive from addr: %v  data: %v\n", remoteAddr, string(result[:n]))
}
