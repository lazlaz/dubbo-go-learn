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

// NewProxy create service proxy.
func NewProxy(invoke protocol.Invoker, callback interface{}, attachments map[string]string) *Proxy {
	return NewProxyWithOptions(invoke, callback, attachments,
		WithProxyImplementFunc(DefaultProxyImplementFunc))
}

// NewProxyWithOptions create service proxy with options.
func NewProxyWithOptions(invoke protocol.Invoker, callback interface{}, attachments map[string]string, opts ...ProxyOption) *Proxy {
	p := &Proxy{
		invoke:      invoke,
		callback:    callback,
		attachments: attachments,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithProxyImplementFunc an option function to setup proxy.ImplementFunc
func WithProxyImplementFunc(f ImplementFunc) ProxyOption {
	return func(p *Proxy) {
		p.implement = f
	}
}

// DefaultProxyImplementFunc the default function for proxy impl
func DefaultProxyImplementFunc(p *Proxy, v common.RPCService) {

}
func (p *Proxy) Implement(v common.RPCService) {
	p.once.Do(func() {
		p.implement(p, v)
		p.rpc = v
	})
}
