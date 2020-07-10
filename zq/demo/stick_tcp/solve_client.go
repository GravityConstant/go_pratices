package main

import (
	"fmt"
	"net"

	"zq/demo/stick_tcp/proto"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.202:30000")
	if err != nil {
		fmt.Println("client dial error: ", err)
		return
	}
	defer conn.Close()

	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		data, err := proto.Encode(msg)
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			return
		}
		conn.Write(data)
	}

}
