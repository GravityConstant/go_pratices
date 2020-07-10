package main

import (
	"errors"
	"log"
	"net/http"
	"net/rpc"
)

type Arith struct{}

type ArithRequest struct {
	A, B int
}

type ArithResponse struct {
	Pro int
	Quo int
	Rem int
}

func (this *Arith) Multiply(req ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

func (this *Arith) Divide(req ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("除数不能为0")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B

	return nil
}

func main() {
	// 注册服务
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
