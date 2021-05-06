package cluster

import "github.com/laz/dubbo-go/protocol"

// Cluster
// Extension - Cluster
type Cluster interface {
	Join(Directory) protocol.Invoker
}
