package main

import (
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	_ "github.com/vlorc/hprose-go-nats"
	"log"
)

func main() {
	client := rpc.NewClient("nats://localhost:4222?topic=test&timeout=1")
	method := &struct{ Hello func(string) (string, error) }{}
	client.UseService(method)
	for i := 0; i < 10; i++ {
		log.Print(method.Hello(fmt.Sprintf("baby(%d)", i)))
	}
}
