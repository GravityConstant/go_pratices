package main

import (
	"fmt"
	"net"
)

// socket_stick/client/main.go

func main() {
	conn, err := net.Dial("tcp", "192.168.1.202:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		conn.Write([]byte(msg))
	}
}
