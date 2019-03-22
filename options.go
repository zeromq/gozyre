// Copyright 2019 The GoZyre Authors. All rights reserved.
// Use of this source code is governed by a MPL-2.0
// license that can be found in the LICENSE file.

package zyre

//#cgo pkg-config: libzyre
//#include<zyre.h>
//void _zyre_set_header(zyre_t *self, const char *group, const char *value) {
//	zyre_set_header(self, group, value);
//}
import "C"

import (
	"fmt"
	"time"
)

// SetHeader - set node header; these are provided to other nodes during
// discovery and come in each ENTER message.
func (z *Node) SetHeader(name string, format string, a ...interface{}) {
	if z.ptr == nil {
		panic("Node.SetHeader: z.ptr is null")
	}
	s := fmt.Sprintf(format, a...)
	C._zyre_set_header(
		z.ptr,
		C.CString(name),
		C.CString(s))
}

func SetHeader(name string, format string, a ...interface{}) Option {
	return func(z *Node) {
		z.SetHeader(name, format, a...)
	}
}

// SetVerbose verbose mode; this tells the node to log all traffic as well as
// all major events.
func (z *Node) SetVerbose() {
	if z.ptr == nil {
		panic("Node.SetVerbose: z.ptr is null")
	}
	C.zyre_set_verbose(z.ptr)
}

func SetVerbose() Option {
	return func(z *Node) {
		z.SetVerbose()
	}
}

// SetPort - Set UDP beacon discovery port; defaults to 5670, this call overrides
// that so you can create independent clusters on the same network, for
// e.g. development vs. production. Has no effect after Start()
func (z *Node) SetPort(port int) {
	if z.ptr == nil {
		panic("Node.SetPort: z.ptr is null")
	}
	C.zyre_set_port(z.ptr, C.int(port))
}

func SetPort(port int) Option {
	return func(z *Node) {
		z.SetPort(port)
	}
}

// SetEvasiveTimeout - Set the peer evasiveness timeout, Default is 5000
// millisecond.  This can be tuned in order to deal with expected network
// conditions and the response time expected by the application. This is tied
// to the beacon interval and rate of messages received.
func (z *Node) SetEvasiveTimeout(interval time.Duration) {
	if z.ptr == nil {
		panic("Node.SetEvasiveTimeout: z.ptr is null")
	}
	C.zyre_set_evasive_timeout(z.ptr, C.int(interval.Nanoseconds()/1000000))
}

func SetEvasiveTimeout(interval time.Duration) Option {
	return func(z *Node) {
		z.SetEvasiveTimeout(interval)
	}
}

// SetExpiredTimeout - Set the peer expiration timeout, default is 30000 milliseconds.
// This can be tuned in order to deal with expected network
// conditions and the response time expected by the application. This is tied
// to the beacon interval and rate of messages received.
func (z *Node) SetExpiredTimeout(interval time.Duration) {
	if z.ptr == nil {
		panic("Node.SetExpiredTimeout: z.ptr is null")
	}
	C.zyre_set_expired_timeout(z.ptr, C.int(interval.Nanoseconds()/1000000))
}

func SetExpiredTimeout(interval time.Duration) Option {
	return func(z *Node) {
		z.SetExpiredTimeout(interval)
	}
}

// SetInterval - Set UDP beacon discovery interval, in milliseconds. Default
// is instant beacon exploration followed by pinging every 1,000 msecs.
func (z *Node) SetInterval(interval time.Duration) {
	if z.ptr == nil {
		panic("Node.SetInterval: z.ptr is null")
	}
	C.zyre_set_interval(z.ptr, C.size_t(interval.Nanoseconds()/1000000))
}

func SetInterval(interval time.Duration) Option {
	return func(z *Node) {
		z.SetInterval(interval)
	}
}

// SetInterval - Set network interface for UDP beacons. If you do not set this,
// CZMQ will choose an interface for you. On boxes with several interfaces you
// should specify which one you want to use, or strange things can happen.
func (z *Node) SetInterface(value string) {
	if z.ptr == nil {
		panic("Node.SetInterface: z.ptr is null")
	}
	C.zyre_set_interface(z.ptr, C.CString(value))
}

func SetInterface(value string) Option {
	return func(z *Node) {
		z.SetInterface(value)
	}
}
