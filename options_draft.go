// +build draft

package zyre

//#cgo pkg-config: libzyre
//#include<zyre.h>
import (
    "C"
)

// SetBeaconPeerPort - Set TCP beacon peer port. The default beacon peer port depends on operating
// but has to be considered as random.
// Use SetBeaconPeerPort() to override the default with a well known value. 
// Very useful to simplify firewall rules (because of randomness of default port).
func (z *Node) SetBeaconPeerPort(port int) {
	if z.ptr == nil {
		panic("Node.SetBeaconPeerPort: z.ptr is null")
	}
	C.zyre_set_beacon_peer_port(z.ptr, C.int(port))
}

func SetBeaconPeerPort(port int) Option {
	return func(z *Node) {
		z.SetBeaconPeerPort(port)
	}
}
