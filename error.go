// Copyright 2018 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package hprose_go_nats

import (
	"errors"
	"net"
)

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

var ErrTimeout = NewTimeout("request timeout")
var ErrNotConnect = errors.New("not connect")
var ErrClosed = errors.New("connect is closed")
var ErrIllegalUrl = errors.New("illegal url")
