package proxy

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/protocol"
	"sync"
)

type (
	// ProxyOption a function to init Proxy with options
	ProxyOption func(p *Proxy)
	// ImplementFunc function for proxy impl of RPCService functions
	ImplementFunc func(p *Proxy, v common.RPCService)
)

// nolint
type Proxy struct {
	rpc         common.RPCService
	invoke      protocol.Invoker
	callback    interface{}
	attachments map[string]string
	implement   ImplementFunc
	once        sync.Once
}
