package main

// minimal zyre app

import (
	"flag"
	"fmt"
	"time"

	zyre "github.com/vyskocilm/gozyre"
)

const GROUP = "GROUP"

func main() {

    println("D: BAF1")
	var name string
	flag.StringVar(&name, "name", "", "node name")
	flag.Parse()

    println("D: BAF2")
	node := zyre.New(name)
	defer node.Destroy()

	err := node.Start()
	if err != nil {
		panic(err)
	}
	err = node.Join(GROUP)
    println("D: BAF3")

	for i := 0; i != 10; i++ {
        println("D: BAF4", i)
		node.Shouts(GROUP, "Hello from %s", name)
		time.Sleep(250 * time.Millisecond)
		msg, err := node.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", msg)
        println("D: BAF5", i)
	}

    println("D: BAF6")
	node.Stop()
}
