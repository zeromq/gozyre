package zyre

import (
	"testing"
)

func TestZyre(t *testing.T) {
	node := New("node")
	defer node.Destroy()
}
