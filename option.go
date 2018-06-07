// Copyright 2018 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package hprose_go_nats

import (
	"crypto/tls"
	"github.com/nats-io/go-nats"
	"net/url"
	"strconv"
	"time"
)

const DEFAULT_QUEUE = 64
const DEFAULT_TIMEOUT = 15
const DEFAULT_DELAY = 45

var queryTable = map[string]func(string, *NatsOption){
	"topic": func(s string, opt *NatsOption) {
		opt.topic = s
	},
	"group": func(s string, opt *NatsOption) {
		opt.group = s
	},
	"id": func(s string, opt *NatsOption) {
		opt.id = s
	},
	"queue": func(s string, opt *NatsOption) {
		opt.queue = int(mustParseInt(s))
	},
	"timeout": func(s string, opt *NatsOption) {
		opt.ttl = mustParseInt(s)
	},
	"delay": func(s string, opt *NatsOption) {
		opt.delay = mustParseInt(s)
	},
}

type NatsOption struct {
	id      string
	topic   string
	group   string
	uri     []string
	queue   int
	ttl     int64
	delay   int64
	options []nats.Option
}

func newDefaultOption(uri ...string) *NatsOption {
	return &NatsOption{
		id:    nats.NewInbox(),
		queue: DEFAULT_QUEUE,
		ttl:   DEFAULT_TIMEOUT,
		delay: DEFAULT_DELAY,
		uri:   uri,
	}
}

func NewOption(option ...func(*NatsOption)) *NatsOption {
	opt := newDefaultOption()
	for _, v := range option {
		v(opt)
	}
	if len(opt.uri) <= 0 {
		opt.uri = append(opt.uri, nats.DefaultURL)
	}
	return opt
}

func setQuery(query url.Values, opt *NatsOption) {
	for k, v := range query {
		if f, ok := queryTable[k]; ok || len(v) > 0 || "" != v[0] {
			f(v[0], opt)
		}
	}
}

func mustParseInt(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if nil != err {
		panic(err)
	}
	return val
}

func Uri(uri ...string) func(*NatsOption) {
	return func(opt *NatsOption) {
		for _, s := range uri {
			u, err := url.Parse(s)
			if nil != err {
				panic(err)
			}
			setQuery(u.Query(), opt)
			if nil != u.User {
				s = u.Scheme + "://" + u.User.String() + "@" + u.Host
			} else {
				s = u.Scheme + "://" + u.Host
			}
			opt.uri = append(opt.uri, s)
		}
	}
}

func Queue(num int) func(*NatsOption) {
	return func(opt *NatsOption) {
		if num > 0 {
			opt.queue = num
		}
	}
}

func RootCAs(ca ...string) func(*NatsOption) {
	return func(opt *NatsOption) {
		if len(ca) > 0 {
			opt.options = append(opt.options, nats.RootCAs(ca...))
		}
	}
}

func ClientCert(certFile, keyFile string) func(*NatsOption) {
	return func(opt *NatsOption) {
		opt.options = append(opt.options, nats.ClientCert(certFile, keyFile))
	}
}

func Secure(tls ...*tls.Config) func(*NatsOption) {
	return func(opt *NatsOption) {
		if len(tls) > 0 {
			opt.options = append(opt.options, nats.Secure(tls...))
		}
	}
}

func Options(o ...nats.Option) func(*NatsOption) {
	return func(opt *NatsOption) {
		if len(o) > 0 {
			opt.options = append(opt.options, o...)
		}
	}
}

func Timeout(t time.Duration) func(*NatsOption) {
	return func(opt *NatsOption) {
		if t > time.Second {
			opt.ttl = int64(t / time.Second)
		}
	}
}

func Delay(t time.Duration) func(*NatsOption) {
	return func(opt *NatsOption) {
		if t > time.Second {
			opt.delay = int64(t / time.Second)
		}
	}
}

func Id(id string) func(*NatsOption) {
	return func(opt *NatsOption) {
		if "" != id {
			opt.id = id
		}
	}
}

func Topic(topic string) func(*NatsOption) {
	return func(opt *NatsOption) {
		if "" != topic {
			opt.topic = topic
		}
	}
}

func Group(group string) func(*NatsOption) {
	return func(opt *NatsOption) {
		opt.group = group
	}
}
