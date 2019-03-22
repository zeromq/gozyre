// Copyright 2019 The GoZyre Authors. All rights reserved.
// Use of this source code is governed by a MPL-2.0
// license that can be found in the LICENSE file.

package zyre

//#cgo pkg-config: libzyre
//#include<zyre.h>
//int _zyre_set_endpoint(zyre_t *self, const char *e) {
//  return zyre_set_endpoint(self, e);
//}
//void _zyre_gossip_bind (zyre_t *self, const char *e) {
//  zyre_gossip_bind(self, e);
//}
//void _zyre_gossip_connect (zyre_t *self, const char *e) {
//  zyre_gossip_connect(self, e);
//}
//int _zyre_whispers(zyre_t *self, const char *peer, const char *s) {
//	return zyre_whispers(self, peer, s);
//}
//int _zyre_shouts(zyre_t *self, const char *group, const char *value) {
//	return zyre_shouts(self, group, value);
//}
//const char *_zlist_firsts(zlist_t *hash) {
//   return (const char*)zlist_first(hash);
//}
//const char *_zlist_nexts(zlist_t *hash) {
//   return (const char*)zlist_next(hash);
//}
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	// ErrStart is returned when node fails to start
	ErrStart = errors.New("zyre_start returned -1")

	// ErrJoin is returned when node fails to join the group
	ErrJoin = errors.New("zyre_join returned -1")

	// ErrLeave is returned when node fails to leave the group
	ErrLeave = errors.New("zyre_leave returned -1")

	// ErrRecvNil is returned when recv got nil pointer
	ErrRecvNil = errors.New("zyre_recv got nil")

	// ErrRecvNilEvent is returned when recv got nil pointer
	ErrRecvNilEvent = errors.New("zyre_recv got nil event")
)

// Node is opaque Golang struct wrapping `zyre_t*`
type Node struct {
	ptr  *C.zyre_t
	uuid string
	name string
}

// New creates a new zyre.Node. Note that until you Start the
// node it is silent and invisible to other nodes on the network.
func New(name string, options ...Option) *Node {
	ptr := C.zyre_new(C.CString(name))
	z := &Node{
		ptr:  ptr,
		uuid: "",
		name: "",
	}
	for _, o := range options {
		o(z)
	}
	return z
}

// Option is a type for setting up the underlying Zyre actor
type Option func(*Node)

// Destroy - destroys a Node node. When you destroy a node, any messages it is
// sending or receiving will be discarded. It frees underlying C memory
func (z *Node) Destroy() {
	if z.ptr == nil {
		return
	}
	C.zyre_destroy(&z.ptr)
	z.ptr = nil
}

// UUID - Return our node UUID string, after successful initialization
func (z *Node) UUID() string {
	if z.ptr == nil {
		panic("Node.UUID: z.ptr is null")
	}
	if z.uuid == "" {
		z.uuid = C.GoString(C.zyre_uuid(z.ptr))
	}
	return z.uuid
}

// Name - return our node name, after successful initialization
func (z *Node) Name() string {
	if z.ptr == nil {
		panic("Node.Name: z.ptr is null")
	}
	if z.name == "" {
		z.name = C.GoString(C.zyre_name(z.ptr))
	}
	return z.name
}

// SetEndpoint - By default, Node binds to an ephemeral TCP port and broadcasts the local
// host name using UDP beaconing. When you call this method, Node will use
// gossip discovery instead of UDP beaconing. You MUST set-up the gossip
// service separately using zyre_gossip_bind() and _connect(). Note that the
// endpoint MUST be valid for both bind and connect operations. You can use
// inproc://, ipc://, or tcp:// transports (for tcp://, use an IP address
// that is meaningful to remote as well as local nodes). Returns error if
// operation zas not succesfull
func (z *Node) SetEndpoint(format string, a ...interface{}) error {
	if z.ptr == nil {
		panic("Node.SetEndpoint: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	rc := C._zyre_set_endpoint(z.ptr, C.CString(s))
	if rc == -1 {
		return fmt.Errorf("Node.SetEndpoint: returned -1")
	}
	return nil
}

// GossipBind - set-up gossip discovery of other nodes. At least one node in the cluster
// must bind to a well-known gossip endpoint, so other nodes can connect to
// it. Note that gossip endpoints are completely distinct from Node node
// endpoints, and should not overlap (they can use the same transport).
func (z *Node) GossipBind(format string, a ...interface{}) {
	if z.ptr == nil {
		panic("Node.GossipBind: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	C._zyre_gossip_bind(z.ptr, C.CString(s))
}

// GossipConnect - Set-up gossip discovery of other nodes. A node may connect to multiple
// other nodes, for redundancy paths. For details of the gossip network
// design, see the CZMQ zgossip class.
func (z *Node) GossipConnect(format string, a ...interface{}) {
	if z.ptr == nil {
		panic("Node.GossipConnect: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	C._zyre_gossip_connect(z.ptr, C.CString(s))
}

// Start - starts a node, after setting header values. When you start a node it
// begins discovery and connection. Returns error if it wasn't
// possible to start the node.
func (z *Node) Start() error {
	if z.ptr == nil {
		panic("Node.Start: z.ptr is null")
	}
	rc := C.zyre_start(z.ptr)
	if rc == -1 {
		return ErrStart
	}
	return nil
}

// Stop node; this signals to other peers that this node will go away.
// This is polite; however you can also just destroy the node without
// stopping it.
func (z *Node) Stop() {
	if z.ptr == nil {
		panic("Node.Stop: z.ptr is null")
	}
	C.zyre_stop(z.ptr)
}

// Join a named group; after joining a group you can send messages to
// the group and all Node nodes in that group will receive them.
func (z *Node) Join(room string) error {
	if z.ptr == nil {
		panic("Node.Join: z.ptr is null")
	}
	rc := C.zyre_join(z.ptr, C.CString(room))
	if rc == -1 {
		return ErrJoin
	}
	return nil
}

// Leave a group
func (z *Node) Leave(room string) error {
	if z.ptr == nil {
		panic("Node.Leave: z.ptr is null")
	}
	rc := C.zyre_leave(z.ptr, C.CString(room))
	if rc == -1 {
		return ErrLeave
	}
	return nil
}

// Recv - Receive next message from network; the message may be a control
// message (Enter, Exit, Join, Leave) or data (Whisper, Shout).
// called must use type switch and type assertions to get the exact type
// Returns error on recv error (unpacking the message, or interrupted)
func (z *Node) Recv() (m interface{}, err error) {
	if z.ptr == nil {
		panic("Node.Recv: z.ptr is null")
	}
	msg := C.zyre_recv(z.ptr)
	if msg == nil {
		err = ErrRecvNil
		return
	}
	defer C.zmsg_destroy(&msg)

	cevent := C.zmsg_popstr(msg)
	if cevent == nil {
		err = ErrRecvNilEvent
		return
	}
	event := C.GoString(cevent)
	defer C.free(unsafe.Pointer(cevent))

	switch event {
	case "ENTER":
		return recvEnter(msg)
	case "EVASIVE":
		return recvEvasive(msg)
	case "EXIT":
		return recvExit(msg)
	case "JOIN":
		return recvJoin(msg)
	case "LEAVE":
		return recvLeave(msg)
	case "WHISPER":
		return recvWhisper(msg)
	case "SHOUT":
		return recvShout(msg)
	case "STOP":
		return recvStop(msg)
	default:
		err = fmt.Errorf("NodeRecv: uknown event '%s'", event)
		return
	}
}

// Whisper - sends byte slice to a single peer specified as UUID string
func (z *Node) Whisper(peer string, data ...[]byte) error {
	if z.ptr == nil {
		panic("Node.Whisper: z.ptr is null")
	}
	msg := C.zmsg_new()
	if msg == nil {
		return fmt.Errorf("Node.Whisper: can't create zmsg_t")
	}
	// we do not defer as zmsg_t will get destroyed ...
	for _, d := range data {
		rc := C.zmsg_addmem(
			msg,
			C.CBytes(d),
			C.size_t(len(data)))
		if rc == -1 {
			C.zmsg_destroy(&msg)
			return fmt.Errorf("Node.Whisper: can't add memory buffer")
		}
	}
	rc := C.zyre_whisper(
		z.ptr,
		C.CString(peer),
		&msg) // .... <- HERE
	if rc == -1 {
		return fmt.Errorf("Node.Whispers failed, returned -1")
	}
	return nil
}

// WhisperString - sends formatted string to a single peer specified as UUID string
func (z *Node) WhisperString(peer string, format string, a ...interface{}) error {
	if z.ptr == nil {
		panic("Node.Whispers: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	rc := C._zyre_whispers(
		z.ptr,
		C.CString(peer),
		C.CString(s))
	if rc == -1 {
		return fmt.Errorf("Node.Whispers failed, returned -1")
	}
	return nil
}

// Shout - sends byte slice to a single peer specified as UUID string
func (z *Node) Shout(group string, data ...[]byte) error {
	if z.ptr == nil {
		panic("Node.Shout: z.ptr is null")
	}
	msg := C.zmsg_new()
	if msg == nil {
		return fmt.Errorf("Node.Shout: can't create zmsg_t")
	}
	// we do not defer as zmsg_t will get destroyed ...
	for _, d := range data {
		rc := C.zmsg_addmem(
			msg,
			C.CBytes(d),
			C.size_t(len(data)))
		if rc == -1 {
			C.zmsg_destroy(&msg)
			return fmt.Errorf("Node.Whisper: can't add memory buffer")
		}
	}
	rc := C.zyre_shout(
		z.ptr,
		C.CString(group),
		&msg) // ... <- HERE
	if rc == -1 {
		return fmt.Errorf("Node.Shout failed, returned -1")
	}
	return nil
}

// ShoutString - Send formatted string to a named group
func (z *Node) ShoutString(group string, format string, a ...interface{}) error {
	if z.ptr == nil {
		panic("Node.Shouts: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	rc := C._zyre_shouts(
		z.ptr,
		C.CString(group),
		C.CString(s))
	if rc == -1 {
		return fmt.Errorf("Node.Shouts failed, returned -1")
	}
	return nil
}

// convert `zlist_t*` to string slice and DESTROY the zlist
func zlistTosliceAndDestroy(list *C.zlist_t) []string {
	if list == nil {
		return []string{}
	}
	defer C.zlist_destroy(&list)
	slice := make([]string, int(C.zlist_size(list)))
	item := C._zlist_firsts(list)
	idx := 0
	for {
		if item == nil {
			break
		}
		slice[idx] = C.GoString(item)
		idx++
		item = C._zlist_nexts(list)
	}
	return slice
}

// Peers - Return zlist of current peer ids.
func (z *Node) Peers() []string {
	if z.ptr == nil {
		panic("Node.Peers: z.ptr is null")
	}
	cpeers := C.zyre_peers(z.ptr)
	return zlistTosliceAndDestroy(cpeers)
}

// PeersByGroup - Return zlist of current peers of this group.
func (z *Node) PeersByGroup(group string) []string {
	if z.ptr == nil {
		panic("Node.PeersByGroup: z.ptr is null")
	}
	cpeers := C.zyre_peers_by_group(z.ptr, C.CString(group))
	return zlistTosliceAndDestroy(cpeers)
}

// PeerGroups Return zlist of current peers of this group.
func (z *Node) PeerGroups() []string {
	if z.ptr == nil {
		panic("Node.PeerGroups: z.ptr is null")
	}
	cpeers := C.zyre_peer_groups(z.ptr)
	return zlistTosliceAndDestroy(cpeers)
}

// PeerAddress - return the endpoint of a connected peer or false if not found
func (z *Node) PeerAddress(peer string) (address string, ok bool) {
	if z.ptr == nil {
		panic("Node.PeerAddress: z.ptr is null")
	}
	caddress := C.zyre_peer_address(z.ptr, C.CString(peer))
	if caddress == nil {
		ok = false
		return
	}
	defer C.free(unsafe.Pointer(caddress))
	address = C.GoString(caddress)
	ok = true
	return
}

// PeerHeaderValue - Return the value of a header of a conected peer.
// Returns ok false if peer or key doesn't exits.
func (z *Node) PeerHeaderValue(peer string, key string) (value string, ok bool) {
	if z.ptr == nil {
		panic("Node.PeerHeaderValue: z.ptr is null")
	}
	cvalue := C.zyre_peer_header_value(z.ptr, C.CString(peer), C.CString(key))
	if cvalue == nil {
		ok = false
		return
	}
	defer C.free(unsafe.Pointer(cvalue))
	ok = true
	value = C.GoString(cvalue)
	return
}
