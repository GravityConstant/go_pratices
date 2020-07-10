package main

import (
	"log"

	"github.com/surgemq/surgemq/service"
)

func main() {
	// Create a new server
	svr := &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
	}
	log.Println("listening tcp on port 1883")
	// Listen and serve connections at localhost:1883
	svr.ListenAndServe("tcp://:1883")
}
