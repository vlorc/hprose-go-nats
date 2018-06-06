package hprose_go_nats

import (
	"github.com/hprose/hprose-golang/rpc"
)

func init() {
	rpc.RegisterClientFactory("nats", newNatsClient)
}
