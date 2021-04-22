package extension

import "github.com/laz/dubbo-go/common/proxy"

var (
	proxyFactories = make(map[string]func(...proxy.Option) proxy.ProxyFactory)
)

// GetProxyFactory finds the ProxyFactory extension with @name
func GetProxyFactory(name string) proxy.ProxyFactory {
	if name == "" {
		name = "default"
	}
	if proxyFactories[name] == nil {
		panic("proxy factory for " + name + " is not existing, make sure you have import the package.")
	}
	return proxyFactories[name]()
}
