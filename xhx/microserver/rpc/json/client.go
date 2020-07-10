package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
)

type Rect struct{}

type Params struct {
	Width, Height int
}

func main() {
	cli, err := jsonrpc.Dial("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	req := Params{10, 5}
	var res int
	if err := cli.Call("Rect.Area", req, &res); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("width:%d, height:%d, area is %d\n", req.Width, req.Height, res)
	}
	if err := cli.Call("Rect.Perimeter", req, &res); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("width:%d, height:%d, area is %d\n", req.Width, req.Height, res)
	}
}
