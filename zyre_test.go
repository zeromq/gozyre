// Copyright 2019 The GoZyre Authors. All rights reserved.
// Use of this source code is governed by a MPL-2.0
// license that can be found in the LICENSE file.

package zyre

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZyre(t *testing.T) {

	assert := assert.New(t)

	node := New(
		"node",
		SetHeader("Service", "name"),
	)
	defer node.Destroy()
	node2 := NewUnique(
		SetPort(5670),
	)
	defer node2.Destroy()
	if testing.Verbose() {
		node.SetVerbose()
		node2.SetVerbose()
	}

	node.SetEvasiveTimeout(5000 * time.Millisecond)
	node.SetExpiredTimeout(30000 * time.Millisecond)
	err := node.Start()
	assert.NoError(err)
	err = node2.Start()
	assert.NoError(err)

	err = node.Join("GROUP")
	assert.NoError(err)
	err = node2.Join("GROUP")
	assert.NoError(err)

	time.Sleep(250 * time.Millisecond)

	for i := 0; i != 5; i++ {
		node.ShoutString("GROUP", "%d#: SHOUT hello from %s", i, "node")
		node.Whisper(node2.UUID(), []byte(fmt.Sprintf("%d#: WHISPER from %s", i, "node")))
	}
	node.Leave("GROUP")

	fmt.Printf("node.Peers=%#v\n", node.Peers())
	fmt.Printf("node2.Peers=%#v\n", node2.Peers())
	a, ok := node.PeerAddress(node2.UUID())
	assert.True(ok)
	fmt.Printf("node.PeerAddress(node2)=%#v\n", a)

	for i := 0; i != 8+5; i++ {
		m, err := node2.Recv()
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
}
