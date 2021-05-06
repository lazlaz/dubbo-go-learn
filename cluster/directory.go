package cluster

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/protocol"
)

// Directory
// Extension - Directory
type Directory interface {
	common.Node
	List(invocation protocol.Invocation) []protocol.Invoker
}
