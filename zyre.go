// Copyright 2019 The GoZyre Authors. All rights reserved.
// Use of this source code is governed by a MPL-2.0
// license that can be found in the LICENSE file.

package zyre

//#cgo pkg-config: libzyre
//#include<zyre.h>
//int _zyre_shouts(zyre_t *self, const char *group, const char *value) {
//	return zyre_shouts(self, group, value);
//}
// void _zyre_print(zyre_t *self) {
//        zmsg_t *msg = zyre_recv(self);
//        zmsg_print(msg);
//        zmsg_destroy (&msg);
// }
import "C"

import (
	"fmt"
    "unsafe"
)

// ZyreEvent - type of even in zyre peer to peer network
type ZyreEvent string
const (
	Enter ZyreEvent = "ENTER"
	Exit  ZyreEvent = "EXIT"
	Join  ZyreEvent = "JOIN"
	Leave  ZyreEvent = "LEAVE"
	Shout ZyreEvent = "SHOUT"
	Evasive ZyreEvent = "EVASIVE"
)

// Zyre - opaque Golang struct wrapping `zyre_t*`
type Zyre struct {
	ptr *C.zyre_t
}

// Message - received zyre message
type Message struct {
	Event ZyreEvent
	Peer  string
	Name  string
	Group string
	Message string
}

// New - creates a new Zyre node. Note that until you start the
// node it is silent and invisible to other nodes on the network.
func New(name string) *Zyre {
	ptr := C.zyre_new(C.CString(name))
	return &Zyre{ptr: ptr}
}

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
	return C.GoString(C.zyre_uuid(z.ptr))
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

// Stop node; this signals to other peers that this node will go away.
// This is polite; however you can also just destroy the node without
// stopping it.
func (z *Zyre) Stop() {
	if z.ptr == nil {
		panic("Zyre.Stop: z.ptr is null")
	}
	C.zyre_stop(z.ptr)
}

// Shouts - Send formatted string to a named group
func (z *Zyre) Shouts(group string, format string, a ...interface{}) error {
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

// Recv - return zyre message or an error
func (z *Zyre) Recv() (m Message, err error) {
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
    event := ZyreEvent(C.GoString(cevent))
    defer C.free(unsafe.Pointer(cevent))

    cpeer := C.zmsg_popstr(msg)
    if cpeer == nil {
        err = fmt.Errorf("Zyre.Recv: got nil peer")
        return
    }
    peer := C.GoString(cpeer)
    defer C.free(unsafe.Pointer(cpeer))

    cname := C.zmsg_popstr(msg)
    if cname == nil {
        err = fmt.Errorf("Zyre.Recv: got nil name")
        return
    }
    name := C.GoString(cname)
    defer C.free(unsafe.Pointer(cname))

    cgroup := C.zmsg_popstr(msg)
    if cgroup == nil {
        err = fmt.Errorf("Zyre.Recv: got nil group")
        return
    }
    group := C.GoString(cgroup)
    defer C.free(unsafe.Pointer(cgroup))

    var message string
    cmessage := C.zmsg_popstr(msg)
    if cmessage == nil {
        message = ""
    } else {
        message = C.GoString(cmessage)
        defer C.free(unsafe.Pointer(cmessage))
    }

    m = Message{
        Event: event,
        Peer: peer,
        Name: name,
        Group: group,
        Message: message,
    }
	return
}
