[![CircleCI](https://circleci.com/gh/gomoni/gozyre.svg?style=svg)](https://circleci.com/gh/gomoni/gozyre)[![license](https://img.shields.io/badge/license-MPL-2.0.svg?style=flat)](https://raw.githubusercontent.com/gomoni/gozyre/master/LICENSE)

# Introduction
A golang interface to the [Zyre v2.0](http://github.com/zeromq/zyre) API. Wraps
most of the API methods to idiomatic Go.

## About Zyre
Zyre is a frameworkd for peer to peer networks written in ZeroMQ

Typical use cases for Zyre are:

 * Local service discovery.
 * Clustering of a set of services on the same Ethernet network.
 * Controlling a network of smart devices (Internet of Things).
 * Multi-user mobile applications (like smart classrooms).

 Check [Zyre Github page](http://github.com/zeromq/zyre)

# Install
## Dependencies
* [libsodium](https://github.com/jedisct1/libsodium)
* [libzmq](https://github.com/zeromq/libzmq)
* [czmq](https://github.com/zeromq/czmq)
* [zyre](https://github.com/zeromq/zyre)

## For GoZyre master
```
import "github.com/gomoni/gozyre"
```

# Example
```go
package main

import (
	"flag"
	"time"

	zyre "github.com/gomoni/gozyre"
)

const (
	Group = "THEROOM"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "node name")
	flag.Parse()

	node := zyre.New(
        name,
        zyre.SetHeader("foo", "bar"),
        )
	defer node.Destroy()
	err := node.Start()
	if err != nil {
		panic(err)
	}
	err = node.Join(Group)
	if err != nil {
		panic(err)
	}
	for i := 0; i != 10; i++ {
		node.Shouts(Group, "Hello from %s", name)
		time.Sleep(250 * time.Millisecond)

		m, err := node.Recv()
		if err != nil {
			fmt.Printf("err=%s\n", err.Error())
		} else {
			switch m.(type) {
			case Shout:
				m := m.(Shout)
				fmt.Printf("%T{Peer:\"%s\", Name:\"%s\", Group:\"%s\", Message: []byte{\"%s\"}}\n", m, m.Peer, m.Name, m.Group, string(m.Message[0]))
			case Whisper:
				m := m.(Whisper)
				fmt.Printf("%T{Peer:\"%s\", Name:\"%s\", Message: []byte{\"%s\"}}\n", m, m.Peer, m.Name, string(m.Message[0]))
			default:
				fmt.Printf("%#+v\n", m)
			}
		}
	}

	node.Stop()
}
```

# Note on panic

`gozyre` panics only when user try to operate on destroyed node

```go
node := zyre.New()
...
node.Destroy()
node.Shout(...) <- panic
```

This is violation of API contract and SHALL not be done.

# License
This project uses the MPL v2 license, see LICENSE.

It uses [Collective Code Construction Contract](https://rfc.zeromq.org/spec:42/C4/) for contributing, see CONTRIBUTING.md.
