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
	"fmt"
	"unsafe"
)

// Zyre - opaque Golang struct wrapping `zyre_t*`
type Zyre struct {
	ptr  *C.zyre_t
	uuid string
	name string
}

// New - creates a new Zyre node. Note that until you start the
// node it is silent and invisible to other nodes on the network.
func New(name string, options ...Option) *Zyre {
	ptr := C.zyre_new(C.CString(name))
	z := &Zyre{
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
type Option func(*Zyre)

// Destroy - destroys a Zyre node. When you destroy a node, any messages it is
// sending or receiving will be discarded. It frees underlying C memory
func (z *Zyre) Destroy() {
	if z.ptr == nil {
		return
	}
	C.zyre_destroy(&z.ptr)
	z.ptr = nil
}

// UUID - Return our node UUID string, after successful initialization
func (z *Zyre) UUID() string {
	if z.ptr == nil {
		panic("Zyre.UUID: z.ptr is null")
	}
	if z.uuid == "" {
		z.uuid = C.GoString(C.zyre_uuid(z.ptr))
	}
	return z.uuid
}

// Name - return our node name, after successful initialization
func (z *Zyre) Name() string {
	if z.ptr == nil {
		panic("Zyre.Name: z.ptr is null")
	}
	if z.name == "" {
		z.name = C.GoString(C.zyre_name(z.ptr))
	}
	return z.name
}

// SetEndpoint - By default, Zyre binds to an ephemeral TCP port and broadcasts the local
// host name using UDP beaconing. When you call this method, Zyre will use
// gossip discovery instead of UDP beaconing. You MUST set-up the gossip
// service separately using zyre_gossip_bind() and _connect(). Note that the
// endpoint MUST be valid for both bind and connect operations. You can use
// inproc://, ipc://, or tcp:// transports (for tcp://, use an IP address
// that is meaningful to remote as well as local nodes). Returns error if
// operation zas not succesfull
func (z *Zyre) SetEndpoint(format string, a ...interface{}) error {
	if z.ptr == nil {
		panic("Zyre.SetEndpoint: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	rc := C._zyre_set_endpoint(z.ptr, C.CString(s))
	if rc == -1 {
		return fmt.Errorf("Zyre.SetEndpoint: returned -1")
	}
	return nil
}

// GossipBind - set-up gossip discovery of other nodes. At least one node in the cluster
// must bind to a well-known gossip endpoint, so other nodes can connect to
// it. Note that gossip endpoints are completely distinct from Zyre node
// endpoints, and should not overlap (they can use the same transport).
func (z *Zyre) GossipBind(format string, a ...interface{}) {
	if z.ptr == nil {
		panic("Zyre.GossipBind: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	C._zyre_gossip_bind(z.ptr, C.CString(s))
}

// GossipConnect - Set-up gossip discovery of other nodes. A node may connect to multiple
// other nodes, for redundancy paths. For details of the gossip network
// design, see the CZMQ zgossip class.
func (z *Zyre) GossipConnect(format string, a ...interface{}) {
	if z.ptr == nil {
		panic("Zyre.GossipConnect: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	C._zyre_gossip_connect(z.ptr, C.CString(s))
}

// Start - starts a node, after setting header values. When you start a node it
// begins discovery and connection. Returns error if it wasn't
// possible to start the node.
func (z *Zyre) Start() error {
	if z.ptr == nil {
		panic("Zyre.Start: z.ptr is null")
	}
	rc := C.zyre_start(z.ptr)
	if rc == -1 {
		return fmt.Errorf("Zyre.Start failed, returned -1")
	}
	return nil
}

// Stop node; this signals to other peers that this node will go away.
// This is polite; however you can also just destroy the node without
// stopping it.
func (z *Zyre) Stop() {
	if z.ptr == nil {
		panic("Zyre.Stop: z.ptr is null")
	}
	C.zyre_stop(z.ptr)
}

// Join a named group; after joining a group you can send messages to
// the group and all Zyre nodes in that group will receive them.
func (z *Zyre) Join(room string) error {
	if z.ptr == nil {
		panic("Zyre.Join: z.ptr is null")
	}
	rc := C.zyre_join(z.ptr, C.CString(room))
	if rc == -1 {
		return fmt.Errorf("Zyre.Join failed, returned -1")
	}
	return nil
}

// Leave a group
func (z *Zyre) Leave(room string) error {
	if z.ptr == nil {
		panic("Zyre.Leave: z.ptr is null")
	}
	rc := C.zyre_leave(z.ptr, C.CString(room))
	if rc == -1 {
		return fmt.Errorf("Zyre.Leave failed, returned -1")
	}
	return nil
}

// Recv - Receive next message from network; the message may be a control
// message (Enter, Exit, Join, Leave) or data (Whisper, Shout).
// called must use type switch and type assertions to get the exact type
// Returns error on recv error (unpacking the message, or interrupted)
func (z *Zyre) Recv() (m interface{}, err error) {
	if z.ptr == nil {
		panic("Zyre.Recv: z.ptr is null")
	}
	msg := C.zyre_recv(z.ptr)
	if msg == nil {
		err = fmt.Errorf("Zyre.Recv: got nil")
		return
	}
	defer C.zmsg_destroy(&msg)

	cevent := C.zmsg_popstr(msg)
	if cevent == nil {
		err = fmt.Errorf("Zyre.Recv: got nil event")
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
		err = fmt.Errorf("ZyreRecv: uknown event '%s'", event)
		return
	}
}

// Whisper - sends byte slice to a single peer specified as UUID string
func (z *Zyre) Whisper(peer string, data ...[]byte) error {
	if z.ptr == nil {
		panic("Zyre.Whisper: z.ptr is null")
	}
	msg := C.zmsg_new()
	if msg == nil {
		return fmt.Errorf("Zyre.Whisper: can't create zmsg_t")
	}
	// we do not defer as zmsg_t will get destroyed ...
	for _, d := range data {
		rc := C.zmsg_addmem(
			msg,
			C.CBytes(d),
			C.size_t(len(data)))
		if rc == -1 {
			C.zmsg_destroy(&msg)
			return fmt.Errorf("Zyre.Whisper: can't add memory buffer")
		}
	}
	rc := C.zyre_whisper(
		z.ptr,
		C.CString(peer),
		&msg) // .... <- HERE
	if rc == -1 {
		return fmt.Errorf("Zyre.Whispers failed, returned -1")
	}
	return nil
}

// WhisperString - sends formatted string to a single peer specified as UUID string
func (z *Zyre) WhisperString(peer string, format string, a ...interface{}) error {
	if z.ptr == nil {
		panic("Zyre.Whispers: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	rc := C._zyre_whispers(
		z.ptr,
		C.CString(peer),
		C.CString(s))
	if rc == -1 {
		return fmt.Errorf("Zyre.Whispers failed, returned -1")
	}
	return nil
}

// Shout - sends byte slice to a single peer specified as UUID string
func (z *Zyre) Shout(group string, data ...[]byte) error {
	if z.ptr == nil {
		panic("Zyre.Shout: z.ptr is null")
	}
	msg := C.zmsg_new()
	if msg == nil {
		return fmt.Errorf("Zyre.Shout: can't create zmsg_t")
	}
	// we do not defer as zmsg_t will get destroyed ...
	for _, d := range data {
		rc := C.zmsg_addmem(
			msg,
			C.CBytes(d),
			C.size_t(len(data)))
		if rc == -1 {
			C.zmsg_destroy(&msg)
			return fmt.Errorf("Zyre.Whisper: can't add memory buffer")
		}
	}
	rc := C.zyre_shout(
		z.ptr,
		C.CString(group),
		&msg) // ... <- HERE
	if rc == -1 {
		return fmt.Errorf("Zyre.Shout failed, returned -1")
	}
	return nil
}

// ShoutString - Send formatted string to a named group
func (z *Zyre) ShoutString(group string, format string, a ...interface{}) error {
	if z.ptr == nil {
		panic("Zyre.Shouts: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	rc := C._zyre_shouts(
		z.ptr,
		C.CString(group),
		C.CString(s))
	if rc == -1 {
		return fmt.Errorf("Zyre.Shouts failed, returned -1")
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
func (z *Zyre) Peers() []string {
	if z.ptr == nil {
		panic("Zyre.Peers: z.ptr is null")
	}
	cpeers := C.zyre_peers(z.ptr)
	return zlistTosliceAndDestroy(cpeers)
}

// PeersByGroup - Return zlist of current peers of this group.
func (z *Zyre) PeersByGroup(group string) []string {
	if z.ptr == nil {
		panic("Zyre.PeersByGroup: z.ptr is null")
	}
	cpeers := C.zyre_peers_by_group(z.ptr, C.CString(group))
	return zlistTosliceAndDestroy(cpeers)
}

// PeerGroups Return zlist of current peers of this group.
func (z *Zyre) PeerGroups() []string {
	if z.ptr == nil {
		panic("Zyre.PeerGroups: z.ptr is null")
	}
	cpeers := C.zyre_peer_groups(z.ptr)
	return zlistTosliceAndDestroy(cpeers)
}

// PeerAddress - return the endpoint of a connected peer or error if not found
func (z *Zyre) PeerAddress(peer string) (address string, err error) {
	if z.ptr == nil {
		panic("Zyre.PeerAddress: z.ptr is null")
	}
	caddress := C.zyre_peer_address(z.ptr, C.CString(peer))
	if caddress == nil {
		err = fmt.Errorf("Zyre.PeerAddress: can't find address of peer %s", peer)
	}
	defer C.free(unsafe.Pointer(caddress))
	address = C.GoString(caddress)
	return
}

// PeerHeaderValue - Return the value of a header of a conected peer.
// Returns ok false if peer or key doesn't exits.
func (z *Zyre) PeerHeaderValue(peer string, key string) (value string, ok bool) {
	if z.ptr == nil {
		panic("Zyre.PeerHeaderValue: z.ptr is null")
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
