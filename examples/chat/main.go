package main

// Example: gozyre distributed chat
// Copyright (c) 2019 The GoZyre Authors

import (
	"bufio"
	"fmt"
	"os"

	zyre "github.com/zeromq/gozyre"
)

func chatActor(pipe chan string, name string) {

	node := zyre.New(
		name,
	)
	node.Start()
	err := node.Join("CHAT")

	zyreChan := make(chan interface{})
	go func() {
		for {
			msg, err := node.Recv()
			if err != nil {
				panic(err)
			}
			zyreChan <- msg
		}
	}()

	if err != nil {
		panic(err)
	}
	for {
		select {
		case msg := <-zyreChan:
			switch msg.(type) {
			case zyre.Join:
				msg := msg.(zyre.Join)
				fmt.Printf("%s has joined the chat\n", msg.Name)
			case zyre.Exit:
				msg := msg.(zyre.Exit)
				fmt.Printf("%s has left the chat\n", msg.Name)
			case zyre.Shout:
				msg := msg.(zyre.Shout)
				fmt.Printf("%s: %s\n", msg.Name, string(msg.Message[0]))
			case zyre.Evasive:
				msg := msg.(zyre.Evasive)
				fmt.Printf("%s is evasive\n", msg.Name)
			}
		case shout := <-pipe:
			if shout == "$TERM" {
				break
			}
			node.ShoutString("CHAT", shout)
		}
	}

	node.Stop()
	node.Destroy()
}

func main() {

	name := os.Args[1]
	pipe := make(chan string)
	go chatActor(pipe, name)

	inp := bufio.NewScanner(os.Stdin)
	for inp.Scan() {
		pipe <- inp.Text()
	}

	pipe <- "$TERM"
}
