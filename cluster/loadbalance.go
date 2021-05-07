package cluster

import "github.com/laz/dubbo-go/protocol"

// LoadBalance
// Extension - LoadBalance
type LoadBalance interface {
	Select([]protocol.Invoker, protocol.Invocation) protocol.Invoker
}
