package main

import (
	"github.com/hprose/hprose-golang/rpc"
	"log"
	_ "worker/hprose-go-nats"
)

func main() {
	client := rpc.NewClient("nats://localhost:4222?topic=cnmb")
	method := &struct{ Hello func(string) (string, error) }{}
	client.UseService(method)
	for i := 0; i < 3000; i++ {
		log.Print(method.Hello("baby"))
	}
}
