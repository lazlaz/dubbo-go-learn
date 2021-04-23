package protocol

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/protocol"
	"sync"
)

var (
	once          sync.Once
	regProtocol   *registryProtocol
	reserveParams = []string{
		"application", "codec", "exchanger", "serialization", "cluster", "connections", "deprecated", "group",
		"loadbalance", "mock", "path", "timeout", "token", "version", "warmup", "weight", "timestamp", "dubbo",
		"release", "interface",
	}
)

type registryProtocol struct {
	invokers []protocol.Invoker
	// Registry Map<RegistryAddress, Registry>
	registries *sync.Map
	// To solve the problem of RMI repeated exposure port conflicts,
	// the services that have been exposed are no longer exposed.
	// providerurl <--> exporter
	bounds                        *sync.Map
	overrideListeners             *sync.Map
	serviceConfigurationListeners *sync.Map

	once sync.Once
}

func init() {
	extension.SetProtocol("registry", GetProtocol)
}

// GetProtocol return the singleton registryProtocol
func GetProtocol() protocol.Protocol {
	//只执行一次
	once.Do(func() {
		regProtocol = newRegistryProtocol()
	})
	return regProtocol
}

// Export provider service to registry center
func (proto *registryProtocol) Export(invoker protocol.Invoker) protocol.Exporter {
	return nil
}

// Destroy registry protocol
func (proto *registryProtocol) Destroy() {

}

// Refer provider service from registry center
func (proto *registryProtocol) Refer(url *common.URL) protocol.Invoker {
	return nil
}
func newRegistryProtocol() *registryProtocol {
	return &registryProtocol{
		registries: &sync.Map{},
		bounds:     &sync.Map{},
	}
}
