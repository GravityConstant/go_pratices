package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"zq/demo/stick_tcp/proto"
)

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("receive from client's message: ", msg)
	}
}

func main() {
	log.Println("listening in port: 30000...")
	listen, err := net.Listen("tcp", ":30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
		}
		go process(conn)
	}
}
