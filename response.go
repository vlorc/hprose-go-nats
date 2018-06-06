package hprose_go_nats

import "github.com/vlorc/timer"

type Request struct {
	seq     Sequence
	wait    chan struct{}
	timeout *timer.Timer
	data    []byte
	err     error
}

func (r *Request) Close() {
	if nil != r.timeout {
		r.timeout.Cancel()
	}
	if nil != r.wait {
		close(r.wait)
	}
}

func (r *Request) Reset() {
	r.wait = nil
	r.timeout = nil
	r.data = nil
	r.err = nil
}

func (r *Request) String() string {
	return r.seq.String()
}

func (r *Request) Id() Sequence {
	return r.seq
}

func (r *Request) Response() ([]byte, error) {
	<-r.wait
	return r.data, r.err
}
