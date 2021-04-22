package proxy

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/protocol"
)

// ProxyFactory interface.
type ProxyFactory interface {
	GetProxy(invoker protocol.Invoker, url *common.URL) *Proxy
	GetAsyncProxy(invoker protocol.Invoker, callBack interface{}, url *common.URL) *Proxy
	GetInvoker(url *common.URL) protocol.Invoker
}

// Option will define a function of handling ProxyFactory
type Option func(ProxyFactory)
