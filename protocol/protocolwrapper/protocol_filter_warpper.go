package protocolwrapper

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/protocol"
)

const (
	// FILTER is protocol key.
	FILTER = "filter"
)

func init() {
	extension.SetProtocol(FILTER, GetProtocol)
}

// ProtocolFilterWrapper
// protocol in url decide who ProtocolFilterWrapper.protocol is
type ProtocolFilterWrapper struct {
	protocol protocol.Protocol
}

// Export service for remote invocation
func (pfw *ProtocolFilterWrapper) Export(invoker protocol.Invoker) protocol.Exporter {
	return nil
}

// Refer a remote service
func (pfw *ProtocolFilterWrapper) Refer(url *common.URL) protocol.Invoker {
	return nil
}

// Destroy will destroy all invoker and exporter.
func (pfw *ProtocolFilterWrapper) Destroy() {
	pfw.protocol.Destroy()
}

// nolint
func GetProtocol() protocol.Protocol {
	return &ProtocolFilterWrapper{}
}
