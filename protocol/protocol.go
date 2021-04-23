package protocol

import (
	"github.com/laz/dubbo-go/common"
	"sync"
)

// Protocol
// Extension - protocol
type Protocol interface {
	// Export service for remote invocation
	Export(invoker Invoker) Exporter
	// Refer a remote service
	Refer(url *common.URL) Invoker
	// Destroy will destroy all invoker and exporter, so it only is called once.
	Destroy()
}

// BaseProtocol is default protocol implement.
type BaseProtocol struct {
	exporterMap *sync.Map
	invokers    []Invoker
}

// Exporter
// wrapping invoker
type Exporter interface {
	// GetInvoker gets invoker.
	GetInvoker() Invoker
	// Unexport exported service.
	Unexport()
}

// NewBaseProtocol creates a new BaseProtocol
func NewBaseProtocol() BaseProtocol {
	return BaseProtocol{
		exporterMap: new(sync.Map),
	}
}
