package zyre

import (
    "fmt"
	"testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestZyre(t *testing.T) {

    assert := assert.New(t)

	node := New("node")
	defer node.Destroy()
    //node.SetVerbose()
    node2 := New("node2")
	defer node2.Destroy()
    //node2.SetVerbose()

    err := node.Start()
    assert.NoError(err)
    err = node2.Start()
    assert.NoError(err)

    err = node.Join("GROUP")
    assert.NoError(err)
    err = node2.Join("GROUP")
    assert.NoError(err)

    time.Sleep(250*time.Millisecond)

    for i := 0; i != 5; i++ {
        node.Shouts("GROUP", "%d#: hello from %s", i, "node")
    }
    node.Leave("GROUP")

    for i := 0; i != 8; i++ {
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
