package main

import (
	"log"
	rpc "worker/hprose-go-nats"
)

func main() {
	server := rpc.NewNatsServer(rpc.Option(rpc.Uri("nats://localhost:4222?topic=cnmb")))
	server.AddFunction("hello", func(msg string) string {
		log.Print("hello: ", msg)
		return "hi bitch!"
	})
	server.Start()
}
