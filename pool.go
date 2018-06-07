// Copyright 2018 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package hprose_go_nats

import (
	"github.com/vlorc/timer"
	"sync"
	"sync/atomic"
)

type Pool struct {
	seq   uint64
	lock  sync.Mutex
	table map[Sequence]*Request
	wheel *timer.TimingWheel
	pool  sync.Pool
}

func NewPool() *Pool {
	return &Pool{
		table: make(map[Sequence]*Request),
		wheel: timer.Default(),
		pool: sync.Pool{
			New: func() interface{} {
				return &Request{}
			},
		},
	}
}

func (p *Pool) Wheel() *timer.TimingWheel {
	p.wheel.Start()
	return p.wheel
}

func (p *Pool) Reset(err error) {
	table := p.table
	tmp := make(map[Sequence]*Request)
	p.lock.Lock()
	p.table = tmp
	p.lock.Unlock()
	for _, v := range table {
		v.err = err
		v.Close()
		p.Put(v)
	}
}

func (p *Pool) Put(req *Request) {
	p.pool.Put(req.Reset())
}

func (p *Pool) Remove(seq Sequence) {
	if req := p.pop(seq); nil != req {
		req.Close()
	}
}

func (p *Pool) pop(seq Sequence) *Request {
	p.lock.Lock()
	res := p.table[seq]
	delete(p.table, seq)
	p.lock.Unlock()
	return res
}

func (p *Pool) get(ttl int64) *Request {
	req := p.pool.Get().(*Request)
	req.seq = Sequence(atomic.AddUint64(&p.seq, 1))
	req.wait = make(chan struct{})
	if ttl > 0 {
		req.timeout = timer.NewTimerTable(func() { p.Push(req.Id(), nil, ErrTimeout) }, ttl)
		p.wheel.Add(req.timeout)
		p.wheel.Start()
	}
	return req
}

func (p *Pool) Get(ttl int64) *Request {
	req := p.get(ttl)
	p.lock.Lock()
	p.table[req.Id()] = req
	p.lock.Unlock()
	return req
}

func (p *Pool) Push(seq Sequence, data []byte, err error) {
	if res := p.pop(seq); nil != res {
		res.data = data
		res.err = err
		res.Close()
	}
}
