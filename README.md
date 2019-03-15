[![CircleCI](https://circleci.com/gh/gomoni/gozyre.svg?style=svg)](https://circleci.com/gh/gomoni/gozyre)

# Introduction
A golang interface to the [Zyre v2.0](http://github.com/zeromq/zyre) API.

# Status

1. Most of methods wrapped 1:1
2. Missing - `zyre_whisper` and `zyre_shout`
3. Should rename `Shouts` to `ShoutString`
4. Examples
5. Options
    - right now there are `Zyre.SetVerbose` and `SetVerbose` methods, maybe keep booth
    - move gossip configuration to extra code
      NewGossip(name, endpoint, bind, connect)

# Install
## Dependencies
* [libsodium](https://github.com/jedisct1/libsodium)
* [libzmq](https://github.com/zeromq/libzmq)
* [czmq](https://github.com/zeromq/czmq)
* [zyre](https://github.com/zeromq/zyre)

## For GoZyre master
```
go get github.com/gomoni/gozyre
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
		//TODO: fix the recv problem!!!
		//C._zyre_print(node.ptr)
	}

	defer node.Destroy()
	node.Stop()
}
```

# License
This project uses the MPL v2 license, see LICENSE
