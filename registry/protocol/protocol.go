package protocol

import (
	"fmt"
	gxset "github.com/dubbogo/gost/container/set"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/protocol"
	"github.com/laz/dubbo-go/protocol/protocolwrapper"
	"github.com/laz/dubbo-go/registry"
	"strings"
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
	proto.once.Do(func() {
		proto.initConfigurationListeners()
	})
	registryUrl := getRegistryUrl(invoker)
	providerUrl := getProviderUrl(invoker)

	overriderUrl := getSubscribedOverrideUrl(providerUrl)
	fmt.Print("overriderUrl %s", overriderUrl)
	var reg registry.Registry
	if regI, loaded := proto.registries.Load(registryUrl.Key()); !loaded {
		reg = getRegistry(registryUrl)
		proto.registries.Store(registryUrl.Key(), reg)
		logger.Infof("Export proto:%p registries address:%p", proto, proto.registries)
	} else {
		reg = regI.(registry.Registry)
	}

	registeredProviderUrl := getUrlToRegistry(providerUrl, registryUrl)
	err := reg.Register(registeredProviderUrl)
	if err != nil {
		logger.Errorf("provider service %v register registry %v error, error message is %s",
			providerUrl.Key(), registryUrl.Key(), err.Error())
		return nil
	}

	key := getCacheKey(providerUrl)
	logger.Infof("The cached exporter keys is %v!", key)
	cachedExporter, loaded := proto.bounds.Load(key)
	if loaded {
		logger.Infof("The exporter has been cached, and will return cached exporter!")
	} else {
		wrappedInvoker := newWrappedInvoker(invoker, providerUrl)
		cachedExporter = extension.GetProtocol(protocolwrapper.FILTER).Export(wrappedInvoker)
		proto.bounds.Store(key, cachedExporter)
		logger.Infof("The exporter has not been cached, and will return a new exporter!")
	}

	return cachedExporter.(protocol.Exporter)
}

type wrappedInvoker struct {
	invoker protocol.Invoker
	protocol.BaseInvoker
}

//包装invoker
func newWrappedInvoker(invoker protocol.Invoker, url *common.URL) *wrappedInvoker {
	return &wrappedInvoker{
		invoker:     invoker,
		BaseInvoker: *protocol.NewBaseInvoker(url),
	}
}
func getCacheKey(url *common.URL) string {
	delKeys := gxset.NewSet("dynamic", "enabled")
	return url.CloneExceptParams(delKeys).String()
}
func getUrlToRegistry(providerUrl *common.URL, registryUrl *common.URL) *common.URL {
	if registryUrl.GetParamBool("simplified", false) {
		return providerUrl.CloneWithParams(reserveParams)
	} else {
		return filterHideKey(providerUrl)
	}
}

// filterHideKey filter the parameters that do not need to be output in url(Starting with .)
func filterHideKey(url *common.URL) *common.URL {
	// be careful params maps in url is map type
	removeSet := gxset.NewSet()
	for k := range url.GetParams() {
		if strings.HasPrefix(k, ".") {
			removeSet.Add(k)
		}
	}
	return url.CloneExceptParams(removeSet)
}

func getRegistry(regUrl *common.URL) registry.Registry {
	reg, err := extension.GetRegistry(regUrl.Protocol, regUrl)
	if err != nil {
		logger.Errorf("Registry can not connect success, program is going to panic.Error message is %s", err.Error())
		panic(err.Error())
	}
	return reg
}
func getSubscribedOverrideUrl(providerUrl *common.URL) *common.URL {
	newUrl := providerUrl.Clone()
	newUrl.Protocol = constant.PROVIDER_PROTOCOL
	newUrl.SetParam(constant.CATEGORY_KEY, constant.CONFIGURATORS_CATEGORY)
	newUrl.SetParam(constant.CHECK_KEY, "false")
	return newUrl
}
func getProviderUrl(invoker protocol.Invoker) *common.URL {
	url := invoker.GetUrl()
	// be careful params maps in url is map type
	return url.SubURL.Clone()
}
func getRegistryUrl(invoker protocol.Invoker) *common.URL {
	// here add * for return a new url
	url := invoker.GetUrl()
	// if the protocol == registry, set protocol the registry value in url.params
	if url.Protocol == constant.REGISTRY_PROTOCOL {
		url.Protocol = url.GetParam(constant.REGISTRY_KEY, "")
	}
	return url
}
func (proto *registryProtocol) initConfigurationListeners() {

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
