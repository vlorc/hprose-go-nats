// Copyright 2018 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package hprose_go_nats

import (
	"github.com/hprose/hprose-golang/rpc"
	"github.com/nats-io/go-nats"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type NatsServer struct {
	rpc.BaseService
	conn        *nats.Conn
	queue       chan *nats.Msg
	signal      chan os.Signal
	uri         string
	opt         *NatsOption
	contextPool sync.Pool
	workerPool  *rpc.WorkerPool
}

func NewServer(opt *NatsOption) rpc.Server {
	server := &NatsServer{
		opt: opt,
		uri: strings.Join(opt.uri, ","),
	}
	server.contextPool.New = func() interface{} {
		return new(rpc.BaseServiceContext)
	}
	server.InitBaseService()
	return server
}

func (ns *NatsServer) worker() {
	for {
		msg, ok := <-ns.queue
		if !ok {
			break
		}
		if nil != msg {
			ns.workerPool.Go(func() {
				ns.handle(msg)
			})
		}
	}
}

func (ns *NatsServer) init() (err error) {
	if nil == ns.conn {
		if ns.conn, err = nats.Connect(ns.uri, ns.opt.options...); nil != err {
			return
		}
	}
	if nil == ns.workerPool {
		ns.workerPool = new(rpc.WorkerPool)
		ns.workerPool.Start()
	}
	if nil == ns.queue {
		ns.queue = make(chan *nats.Msg, ns.opt.queue)
	}
	if _, err = ns.conn.ChanQueueSubscribe(ns.opt.topic, ns.opt.group, ns.queue); nil != err {
		return
	}
	go ns.worker()
	if nil == ns.signal {
		ns.signal = make(chan os.Signal, 1)
		signal.Notify(ns.signal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	}
	return
}
func (ns *NatsServer) handle(msg *nats.Msg) {
	ctx := ns.contextPool.Get().(*rpc.BaseServiceContext)
	defer ns.contextPool.Put(ctx)

	ctx.InitServiceContext(ns)
	data := ns.BaseService.Handle(msg.Data, ctx)
	ns.conn.Publish(msg.Reply, data)
}

func (ns *NatsServer) URI() string {
	return ns.uri
}

func (ns *NatsServer) Handle() error {
	return nil
}

func (ns *NatsServer) Close() {
	if nil != ns.signal {
		signal.Stop(ns.signal)
		ns.signal = nil
	}
	if nil != ns.queue {
		close(ns.queue)
		ns.queue = nil
	}
	if nil != ns.conn {
		ns.conn.Close()
		ns.conn = nil
	}
	if nil != ns.workerPool {
		ns.workerPool.Stop()
		ns.workerPool = nil
	}
}

func (ns *NatsServer) Start() (err error) {
	if err = ns.init(); nil != err {
		return
	}
	<-ns.signal
	ns.Close()
	return nil
}

func (ns *NatsServer) Restart() {
	ns.signal <- syscall.SIGHUP
}

func (ns *NatsServer) Stop() {
	ns.signal <- syscall.SIGQUIT
}
