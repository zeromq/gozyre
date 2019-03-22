// Copyright 2019 The GoZyre Authors. All rights reserved.
// Use of this source code is governed by a MPL-2.0
// license that can be found in the LICENSE file.

package zyre

//#cgo pkg-config: libzyre
//#include<zyre.h>
//const char *zhash_firsts(zhash_t *hash) {
//   return (const char*)zhash_first(hash);
//}
//const char *zhash_nexts(zhash_t *hash) {
//   return (const char*)zhash_next(hash);
//}
//const char *zhash_cursors(zhash_t *hash) {
//   return (const char*)zhash_cursor(hash);
//}
import "C"

import (
	"fmt"
	"unsafe"
)

// Enter - new peer has entered the network
type Enter struct {
	Peer     string
	Name     string
	Headers  map[string]string
	Endpoint string
}

// Evasive - peer is being evasive (quiet for too long)
type Evasive struct {
	Peer string
	Name string
}

// Exit - peer has left the network
type Exit struct {
	Peer string
	Name string
}

// Join - peer has joined a specific group
type Join struct {
	Peer  string
	Name  string
	Group string
}

// Leave - peer has left a specific group
type Leave struct {
	Peer  string
	Name  string
	Group string
}

// Whisper -  peer has sent this node a message
type Whisper struct {
	Peer    string
	Name    string
	Message [][]byte
}

// Shout -  a peer has sent one of our groups a message
type Shout struct {
	Peer    string
	Name    string
	Group   string
	Message [][]byte
}

// Stop - peer was stopped
type Stop struct {
	Peer string
	Name string
}


func recvEnter(msg *C.zmsg_t) (m Enter, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: ENTER got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: ENTER got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	cheaders := C.zmsg_pop(msg)
	if cheaders == nil {
		err = fmt.Errorf("Zyre.Recv: ENTER got nil headers")
		return
	}
	defer C.free(unsafe.Pointer(cheaders))
	chash := C.zhash_unpack(cheaders)
	if chash == nil {
		err = fmt.Errorf("Zyre.Recv: ENTER headers unpack failed")
		return
	}
	defer C.zhash_destroy(&chash)
	headers := make(map[string]string)
	value := C.zhash_firsts(chash)
	for value != nil {
		key := C.zhash_cursors(chash)
		headers[C.GoString(key)] = C.GoString(value)
		value = C.zhash_nexts(chash)
	}

	cip := C.zmsg_popstr(msg)
	if cip == nil {
		err = fmt.Errorf("Zyre.Recv: ENTER got nil ip:port")
		return
	}
	ip := C.GoString(cip)
	defer C.free(unsafe.Pointer(cip))

	m = Enter{
		Peer:     peer,
		Name:     name,
		Headers:  headers,
		Endpoint: ip,
	}
	return
}

func recvEvasive(msg *C.zmsg_t) (m Evasive, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: EVASIVE got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: EVASIVE got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	m = Evasive{
		Peer: peer,
		Name: name,
	}
	return
}

func recvExit(msg *C.zmsg_t) (m Exit, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: EXIT got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: EXIT got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	m = Exit{
		Peer: peer,
		Name: name,
	}
	return
}

func recvJoin(msg *C.zmsg_t) (m Join, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: JOIN got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: JOIN got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	cgroup := C.zmsg_popstr(msg)
	if cgroup == nil {
		err = fmt.Errorf("Zyre.Recv: JOIN got nil group")
		return
	}
	group := C.GoString(cgroup)
	defer C.free(unsafe.Pointer(cgroup))

	m = Join{
		Peer:  peer,
		Name:  name,
		Group: group,
	}
	return
}

func recvLeave(msg *C.zmsg_t) (m Leave, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: LEAVE got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: LEAVE got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	cgroup := C.zmsg_popstr(msg)
	if cgroup == nil {
		err = fmt.Errorf("Zyre.Recv: LEAVE got nil group")
		return
	}
	group := C.GoString(cgroup)
	defer C.free(unsafe.Pointer(cgroup))

	m = Leave{
		Peer:  peer,
		Name:  name,
		Group: group,
	}
	return
}

func recvWhisper(msg *C.zmsg_t) (m Whisper, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: WHISPER got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: WHISPER got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	message := make([][]byte, C.zmsg_size(msg))
	i := 0
	for {
		p := C.zmsg_pop(msg)
		if p == nil {
			break
		}
		message[i] = C.GoBytes(
			unsafe.Pointer(C.zframe_data(p)),
			C.int(C.zframe_size(p)),
		)
	}

	m = Whisper{
		Peer:    peer,
		Name:    name,
		Message: message,
	}
	return
}

func recvShout(msg *C.zmsg_t) (m Shout, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: SHOUT got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: SHOUT got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	cgroup := C.zmsg_popstr(msg)
	if cgroup == nil {
		err = fmt.Errorf("Zyre.Recv: SHOUT got nil group")
		return
	}
	group := C.GoString(cgroup)
	defer C.free(unsafe.Pointer(cgroup))

	message := make([][]byte, C.zmsg_size(msg))
	i := 0
	for {
		p := C.zmsg_pop(msg)
		if p == nil {
			break
		}
		message[i] = C.GoBytes(
			unsafe.Pointer(C.zframe_data(p)),
			C.int(C.zframe_size(p)),
		)
	}

	m = Shout{
		Peer:    peer,
		Name:    name,
		Group:   group,
		Message: message,
	}
	return
}

func recvStop(msg *C.zmsg_t) (m Stop, err error) {
	cpeer := C.zmsg_popstr(msg)
	if cpeer == nil {
		err = fmt.Errorf("Zyre.Recv: STOP got nil peer")
		return
	}
	peer := C.GoString(cpeer)
	defer C.free(unsafe.Pointer(cpeer))

	cname := C.zmsg_popstr(msg)
	if cname == nil {
		err = fmt.Errorf("Zyre.Recv: STOP got nil name")
		return
	}
	name := C.GoString(cname)
	defer C.free(unsafe.Pointer(cname))

	m = Stop{
		Peer:    peer,
		Name:    name,
	}
	return
}
