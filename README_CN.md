<p align="center"><img src="http://hprose.com/banner.@2x.png" alt="Hprose" title="Hprose" width="650" height="200" /></p>

# [Hprose for nats](https://github.com/vlorc/hprose-go-nats)
[English](https://github.com/vlorc/hprose-go-nats/blob/master/README.md)

[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![codebeat badge](https://codebeat.co/badges/c41b426c-4121-4dc8-99c2-f1b60574be64)](https://codebeat.co/projects/github-com-vlorc-hprose-go-nats-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/vlorc/hprose-go-nats)](https://goreportcard.com/report/github.com/vlorc/hprose-go-nats)
[![GoDoc](https://godoc.org/github.com/vlorc/hprose-go-nats?status.svg)](https://godoc.org/github.com/vlorc/hprose-go-nats)
[![Build Status](https://travis-ci.org/vlorc/hprose-go-nats.svg?branch=master)](https://travis-ci.org/vlorc/hprose-go-nats?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/vlorc/hprose-go-nats/badge.svg?branch=master)](https://coveralls.io/github/vlorc/gioc?branch=master)

基于NATs消息队列的Hprose

## 特性
+ 超时
+ 惰性加载
+ 失败切换
+ 负载均衡

## 安装
	go get github.com/vlorc/hprose-go-nats

## 许可证
这个项目是在Apache许可证下进行的。请参阅完整许可证文本的许可证文件。

## 实例
### 客户端
```golang
client := rpc.NewClient("nats://localhost:4222?topic=test&timeout=1")
method := &struct{ Hello func(string) (string, error) }{}
client.UseService(method)
for i := 0; i < 10; i++ {
	log.Print(method.Hello(fmt.Sprintf("baby(%d)",i)))
}
```

### 服务端
```golang
server := rpc.NewServer(rpc.NewOption(rpc.Uri("nats://localhost:4222?topic=test&group=balancer")))
server.AddFunction("hello", func(msg string) string {
	log.Print("hello: ", msg)
	return "hi bitch!"
})
server.Start()
```