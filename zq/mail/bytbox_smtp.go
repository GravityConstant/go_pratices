package main

import (
	"log"

	"github.com/bytbox/go-pop3"
)

var (
	address = "pop.qiye.163.com:110"
	user    = "php01@2009400.cn"
	pass    = "ZQxx400@"
)

func main() {
	cmdSendMail()
}

func cmdSendMail() {
	client, err := pop3.Dial(address)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	defer func() {
		client.Quit()
	}()

	if err = client.User(user); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if err = client.Pass(pass); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	var content string

	if content, err = client.Retr(1); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	log.Printf("Content:\n%s\n", content)
}
