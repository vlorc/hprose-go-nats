package main

import (
	rpc "github.com/vlorc/hprose-go-nats"
	"log"
)

func main() {
	server := rpc.NewServer(rpc.NewOption(rpc.Uri("nats://localhost:4222?topic=test&group=balancer")))
	server.AddFunction("hello", func(msg string) string {
		log.Print("hello: ", msg)
		return "hi bitch!"
	})
	server.Start()
}
