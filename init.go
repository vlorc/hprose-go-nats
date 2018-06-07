// Copyright 2018 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package hprose_go_nats

import (
	"github.com/hprose/hprose-golang/rpc"
)

func init() {
	rpc.RegisterClientFactory("nats", func(uri ...string) rpc.Client {
		return NewClient(newDefaultOption(uri...))
	})
}
