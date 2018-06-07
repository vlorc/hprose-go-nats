// Copyright 2018 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package hprose_go_nats

import (
	"errors"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/nats-io/go-nats"
	"github.com/vlorc/timer"
	"net/url"
	"sync"
	"sync/atomic"
)

type NatsClient struct {
	rpc.BaseClient
	conn  *nats.Conn
	url   *url.URL
	uri   string
	opt   *NatsOption
	queue chan *nats.Msg
	pool  *Pool
	lock  sync.Mutex
	send  atomic.Value
}

func NewClient(opt *NatsOption) rpc.Client {
	client := &NatsClient{
		queue: make(chan *nats.Msg, opt.queue),
		pool:  NewPool(),
		opt:   opt,
	}
	client.InitBaseClient()
	client.SetURIList(opt.uri)
	client.send.Store(sendErrNotConnect)
	client.BaseClient.SendAndReceive = client.SendAndReceive
	go client.worker()
	return client
}

func (nc *NatsClient) worker() {
	for {
		msg, ok := <-nc.queue
		if !ok {
			break
		}
		if nil != msg {
			nc.pool.Push(NewSequence(msg.Subject), msg.Data, nil)
		}
	}
}

func (nc *NatsClient) sendRequest(data []byte, ctx *rpc.ClientContext) ([]byte, error) {
	req := nc.pool.Get(nc.opt.ttl)
	reply := nc.opt.id + req.String()
	if err := nc.conn.PublishRequest(nc.opt.topic, reply, data); nil != err {
		nc.pool.Remove(req.Id())
		return nil, err
	}
	return req.Response()
}

func (nc *NatsClient) ready() *NatsClient {
	if nc.url != nc.BaseClient.URL() {
		nc.lock.Lock()
		if nc.url != nc.BaseClient.URL() {
			nc.init()
			nc.url = nc.BaseClient.URL()
		}
		nc.lock.Unlock()
	}
	return nc
}

func (nc *NatsClient) SendAndReceive(data []byte, ctx *rpc.ClientContext) ([]byte, error) {
	return nc.ready().send.Load().(func([]byte, *rpc.ClientContext) ([]byte, error))(data, ctx)
}

func (nc *NatsClient) shutdown() {
	if nil != nc.conn {
		conn := nc.conn
		nc.conn = nil
		nc.pool.Wheel().Add(timer.NewTimerTable(func() { conn.Close() }, nc.opt.delay))
	}
}

func (nc *NatsClient) Close() {
	nc.send.Store(sendErrClosed)
	nc.shutdown()
}

func (nc *NatsClient) init() {
	if nil == nc.BaseClient.URL() {
		nc.send.Store(sendErrIllegalUrl)
		return
	}
	if nc.uri == nc.BaseClient.URI() {
		return
	}
	nc.connect()
}

func (nc *NatsClient) connect() {
	var err error
	defer func() {
		if it := recover(); nil != it {
			err = errors.New(fmt.Sprint(it))
		}
		if nil != err {
			nc.send.Store(func([]byte, *rpc.ClientContext) ([]byte, error) {
				return nil, err
			})
		}
	}()
	nc.shutdown()
	nc.opt.uri = nc.opt.uri[:0]
	Uri(nc.BaseClient.URI())(nc.opt)
	if nc.conn, err = nats.Connect(nc.opt.uri[0], nc.opt.options...); nil != err {
		return
	}
	_, err = nc.conn.ChanSubscribe(nc.opt.id+".*", nc.queue)
	if nil != err {
		return
	}
	nc.uri = nc.BaseClient.URI()
	nc.send.Store(nc.sendRequest)
}

func sendErrNotConnect([]byte, *rpc.ClientContext) ([]byte, error) {
	return nil, ErrNotConnect
}

func sendErrClosed([]byte, *rpc.ClientContext) ([]byte, error) {
	return nil, ErrClosed
}

func sendErrIllegalUrl([]byte, *rpc.ClientContext) ([]byte, error) {
	return nil, ErrIllegalUrl
}
