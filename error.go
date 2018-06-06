package hprose_go_nats

import "net"

type timeout func() string

func (timeout) Timeout() bool {
	return true
}

func (timeout) Temporary() bool {
	return false
}

func (t timeout) Error() string {
	return t()
}

func NewTimeout(s string) net.Error {
	return timeout(func() string {
		return s
	})
}

func NewTimeoutFunc(f func() string) net.Error {
	return timeout(f)
}
