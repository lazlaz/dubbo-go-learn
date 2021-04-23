package proxy_factory

import (
	"context"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/common/proxy"
	"github.com/laz/dubbo-go/protocol"
)

// DefaultProxyFactory is the default proxy factory
type DefaultProxyFactory struct {
	//delegate ProxyFactory
}

// ProxyInvoker is a invoker struct
type ProxyInvoker struct {
	protocol.BaseInvoker
}

func init() {
	extension.SetProxyFactory("default", NewDefaultProxyFactory)
}

// NewDefaultProxyFactory returns a proxy factory instance
func NewDefaultProxyFactory(_ ...proxy.Option) proxy.ProxyFactory {
	return &DefaultProxyFactory{}
}

// GetProxy gets a proxy
func (factory *DefaultProxyFactory) GetProxy(invoker protocol.Invoker, url *common.URL) *proxy.Proxy {
	return factory.GetAsyncProxy(invoker, nil, url)
}

// GetAsyncProxy gets a async proxy
func (factory *DefaultProxyFactory) GetAsyncProxy(invoker protocol.Invoker, callBack interface{}, url *common.URL) *proxy.Proxy {
	//create proxy
	attachments := map[string]string{}
	attachments[constant.ASYNC_KEY] = url.GetParam(constant.ASYNC_KEY, "false")
	return proxy.NewProxy(invoker, callBack, attachments)
}

// GetInvoker gets a invoker
func (factory *DefaultProxyFactory) GetInvoker(url *common.URL) protocol.Invoker {
	return &ProxyInvoker{
		BaseInvoker: *protocol.NewBaseInvoker(url),
	}
}

// Invoke is used to call service method by invocation
func (pi *ProxyInvoker) Invoke(ctx context.Context, invocation protocol.Invocation) protocol.Result {
	return nil
}
