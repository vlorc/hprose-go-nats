package hprose_go_nats

import (
	"github.com/nats-io/go-nats"
	"net/url"
	"strconv"
	"time"
)

type NatsServerOption struct {
	id      string
	topic   string
	group   string
	uri     []string
	queue   int
	timeout time.Duration
	options []nats.Option
}

func Option(option ...func(*NatsServerOption)) *NatsServerOption {
	opt := &NatsServerOption{
		id:    nats.NewInbox(),
		queue: DEFAULT_QUEUE,
	}
	for _, v := range option {
		v(opt)
	}
	return opt
}

func setQuery(val url.Values, opt *NatsServerOption) {
	if len(val) <= 0 {
		return
	}
	if "" != val.Get("topic") {
		opt.topic = val.Get("topic")
	}
	if "" != val.Get("group") {
		opt.group = val.Get("topic")
	}
	if "" != val.Get("id") {
		opt.id = val.Get("id")
	}
	if "" != val.Get("queue") {
		if queue, _ := strconv.ParseInt(val.Get("queue"), 10, 64); queue > 0 {
			opt.queue = int(queue)
		}
	}
}

func Uri(uri ...string) func(*NatsServerOption) {
	return func(opt *NatsServerOption) {
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

func Queue(num int) func(*NatsServerOption) {
	return func(opt *NatsServerOption) {
		if num > 0 {
			opt.queue = num
		}
	}
}

func NatsOption(o ...nats.Option) func(*NatsServerOption) {
	return func(opt *NatsServerOption) {
		if len(o) > 0 {
			opt.options = append(opt.options, o...)
		}
	}
}

func Timeout(t time.Duration) func(*NatsServerOption) {
	return func(opt *NatsServerOption) {
		if t > 0 {
			opt.timeout = t
		}
	}
}

func Id(id string) func(*NatsServerOption) {
	return func(opt *NatsServerOption) {
		if "" != id {
			opt.id = id
		}
	}
}

func Topic(topic string) func(*NatsServerOption) {
	return func(opt *NatsServerOption) {
		if "" != topic {
			opt.topic = topic
		}
	}
}

func Group(group string) func(*NatsServerOption) {
	return func(opt *NatsServerOption) {
		opt.group = group
	}
}
