package main

import (
	"fmt"
	"net"
)

func main() {
	socket, err := net.DialUDP("udp", nil,
		&net.UDPAddr{
			IP:   net.IPv4(192, 168, 1, 202),
			Port: 12000,
		})
	if err != nil {
		fmt.Println("connect to server failed. err: ", err)
		return
	}
	defer socket.Close()
	sendData := []byte("Hello server")
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("send data error. err:", err)
		return
	}
	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("receive data error. err:", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}
