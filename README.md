<p align="center"><img src="http://hprose.com/banner.@2x.png" alt="Hprose" title="Hprose" width="650" height="200" /></p>

# [Hprose for nats](https://github.com/vlorc/hprose-go-nats)
[简体中文](https://github.com/vlorc/hprose-go-nats/blob/master/README_CN.md)

[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![codebeat badge](https://codebeat.co/badges/c41b426c-4121-4dc8-99c2-f1b60574be64)](https://codebeat.co/projects/github-com-vlorc-hprose-go-nats-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/vlorc/hprose-go-nats)](https://goreportcard.com/report/github.com/vlorc/hprose-go-nats)
[![GoDoc](https://godoc.org/github.com/vlorc/hprose-go-nats?status.svg)](https://godoc.org/github.com/vlorc/hprose-go-nats)
[![Build Status](https://travis-ci.org/vlorc/hprose-go-nats.svg?branch=master)](https://travis-ci.org/vlorc/hprose-go-nats?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/vlorc/hprose-go-nats/badge.svg?branch=master)](https://coveralls.io/github/vlorc/gioc?branch=master)

Hprose based on NATs message queue

## Features
+ timeout
+ lazy load
+ failround
+ load balancing

## Installing
	go get github.com/vlorc/hprose-go-nats

## License
This project is under the apache License. See the LICENSE file for the full license text.

## Examples
### Client
```golang
client := rpc.NewClient("nats://localhost:4222?topic=test&timeout=1")
method := &struct{ Hello func(string) (string, error) }{}
client.UseService(method)
for i := 0; i < 10; i++ {
	log.Print(method.Hello(fmt.Sprintf("baby(%d)",i)))
}
```

### Server
```golang
server := rpc.NewServer(rpc.NewOption(rpc.Uri("nats://localhost:4222?topic=test&group=balancer")))
server.AddFunction("hello", func(msg string) string {
	log.Print("hello: ", msg)
	return "hi bitch!"
})
server.Start()
```
