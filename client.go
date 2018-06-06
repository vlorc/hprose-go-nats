package hprose_go_nats

import (
	"errors"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/nats-io/go-nats"
	"github.com/vlorc/timer"
	"net/url"
	"sync"
	"time"
)

const DEFAULT_QUEUE = 64

type NatsClient struct {
	rpc.BaseClient
	conn  *nats.Conn
	url   *url.URL
	uri   string
	opt   *NatsServerOption
	queue chan *nats.Msg
	pool  *Pool
	lock  sync.Mutex
	send  func([]byte, *rpc.ClientContext) ([]byte, error)
}

func newNatsClient(uri ...string) rpc.Client {
	return NewNatsClient(&NatsServerOption{id: nats.NewInbox(), uri: uri})
}

func NewNatsClient(opt *NatsServerOption) rpc.Client {
	client := &NatsClient{
		queue: make(chan *nats.Msg, opt.queue),
		pool:  NewPool(),
		opt:   opt,
	}
	client.InitBaseClient()
	client.SetURIList(opt.uri)
	opt.uri = nil
	client.send = client.sendRequest
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
	ttl := int64(nc.Timeout() / time.Second)
	req := nc.pool.Get(ttl)
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
	return nc.ready().send(data, ctx)
}

func (nc *NatsClient) shutdown(ttl int64) {
	if nil != nc.conn {
		conn := nc.conn
		nc.conn = nil
		nc.pool.Wheel().Add(timer.NewTimerTable(func() {
			conn.Close()
		}, ttl))
	}
}

func (nc *NatsClient) Close() {
	err := errors.New("client is closed")
	nc.send = func(data []byte, ctx *rpc.ClientContext) ([]byte, error) {
		return nil, err
	}
	nc.shutdown(60)
}

func (nc *NatsClient) init() {
	if nc.uri == nc.BaseClient.URI() {
		return
	}
	var err error
	defer func() {
		if it := recover(); nil != it {
			err = errors.New(fmt.Sprint(it))
		}
		if nil != err {
			nc.send = func([]byte, *rpc.ClientContext) ([]byte, error) {
				return nil, err
			}
		}
	}()
	nc.shutdown(60)
	nc.opt.uri = nil
	Uri(nc.BaseClient.URI())(nc.opt)
	if nc.conn, err = nats.Connect(nc.opt.uri[0], nc.opt.options...); nil != err {
		return
	}
	_, err = nc.conn.ChanSubscribe(nc.opt.id+".*", nc.queue)
	if nil != err {
		return
	}
	nc.uri = nc.BaseClient.URI()
}
