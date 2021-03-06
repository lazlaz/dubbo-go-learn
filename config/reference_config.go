package config

import (
	"context"
	"fmt"
	gxstrings "github.com/dubbogo/gost/strings"
	"github.com/laz/dubbo-go/cluster/directory"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/common/proxy"
	"github.com/laz/dubbo-go/protocol"
	"net/url"
	"strconv"
	"time"
)

// ReferenceConfig is the configuration of service consumer
type ReferenceConfig struct {
	context        context.Context
	pxy            *proxy.Proxy
	id             string
	InterfaceName  string            `required:"true"  yaml:"interface"  json:"interface,omitempty" property:"interface"`
	Check          *bool             `yaml:"check"  json:"check,omitempty" property:"check"`
	Url            string            `yaml:"url"  json:"url,omitempty" property:"url"`
	Filter         string            `yaml:"filter" json:"filter,omitempty" property:"filter"`
	Protocol       string            `default:"dubbo"  yaml:"protocol"  json:"protocol,omitempty" property:"protocol"`
	Registry       string            `yaml:"registry"  json:"registry,omitempty"  property:"registry"`
	Cluster        string            `yaml:"cluster"  json:"cluster,omitempty" property:"cluster"`
	Loadbalance    string            `yaml:"loadbalance"  json:"loadbalance,omitempty" property:"loadbalance"`
	Retries        string            `yaml:"retries"  json:"retries,omitempty" property:"retries"`
	Group          string            `yaml:"group"  json:"group,omitempty" property:"group"`
	Version        string            `yaml:"version"  json:"version,omitempty" property:"version"`
	ProvideBy      string            `yaml:"provide_by"  json:"provide_by,omitempty" property:"provide_by"`
	Methods        []*MethodConfig   `yaml:"methods"  json:"methods,omitempty" property:"methods"`
	Async          bool              `yaml:"async"  json:"async,omitempty" property:"async"`
	Params         map[string]string `yaml:"params"  json:"params,omitempty" property:"params"`
	invoker        protocol.Invoker
	urls           []*common.URL
	Generic        bool   `yaml:"generic"  json:"generic,omitempty" property:"generic"`
	Sticky         bool   `yaml:"sticky"   json:"sticky,omitempty" property:"sticky"`
	RequestTimeout string `yaml:"timeout"  json:"timeout,omitempty" property:"timeout"`
	ForceTag       bool   `yaml:"force.tag"  json:"force.tag,omitempty" property:"force.tag"`
}

func (c *ReferenceConfig) getUrlMap() url.Values {
	urlMap := url.Values{}
	//first set user params
	for k, v := range c.Params {
		urlMap.Set(k, v)
	}
	urlMap.Set(constant.INTERFACE_KEY, c.InterfaceName)
	urlMap.Set(constant.TIMESTAMP_KEY, strconv.FormatInt(time.Now().Unix(), 10))
	urlMap.Set(constant.CLUSTER_KEY, c.Cluster)
	urlMap.Set(constant.LOADBALANCE_KEY, c.Loadbalance)
	urlMap.Set(constant.RETRIES_KEY, c.Retries)
	urlMap.Set(constant.GROUP_KEY, c.Group)
	urlMap.Set(constant.VERSION_KEY, c.Version)
	urlMap.Set(constant.GENERIC_KEY, strconv.FormatBool(c.Generic))
	urlMap.Set(constant.ROLE_KEY, strconv.Itoa(common.CONSUMER))
	urlMap.Set(constant.PROVIDER_BY, c.ProvideBy)

	urlMap.Set(constant.RELEASE_KEY, "dubbo-golang-"+constant.Version)
	urlMap.Set(constant.SIDE_KEY, (common.RoleType(common.CONSUMER)).Role())

	if len(c.RequestTimeout) != 0 {
		urlMap.Set(constant.TIMEOUT_KEY, c.RequestTimeout)
	}
	//getty invoke async or sync
	urlMap.Set(constant.ASYNC_KEY, strconv.FormatBool(c.Async))
	urlMap.Set(constant.STICKY_KEY, strconv.FormatBool(c.Sticky))

	//application info
	urlMap.Set(constant.APPLICATION_KEY, consumerConfig.ApplicationConfig.Name)
	urlMap.Set(constant.ORGANIZATION_KEY, consumerConfig.ApplicationConfig.Organization)
	urlMap.Set(constant.NAME_KEY, consumerConfig.ApplicationConfig.Name)
	urlMap.Set(constant.MODULE_KEY, consumerConfig.ApplicationConfig.Module)
	urlMap.Set(constant.APP_VERSION_KEY, consumerConfig.ApplicationConfig.Version)
	urlMap.Set(constant.OWNER_KEY, consumerConfig.ApplicationConfig.Owner)
	urlMap.Set(constant.ENVIRONMENT_KEY, consumerConfig.ApplicationConfig.Environment)

	//filter
	var defaultReferenceFilter = constant.DEFAULT_REFERENCE_FILTERS
	if c.Generic {
		defaultReferenceFilter = constant.GENERIC_REFERENCE_FILTERS + "," + defaultReferenceFilter
	}
	urlMap.Set(constant.REFERENCE_FILTER_KEY, mergeValue(consumerConfig.Filter, c.Filter, defaultReferenceFilter))

	for _, v := range c.Methods {
		urlMap.Set("methods."+v.Name+"."+constant.LOADBALANCE_KEY, v.LoadBalance)
		urlMap.Set("methods."+v.Name+"."+constant.RETRIES_KEY, v.Retries)
		urlMap.Set("methods."+v.Name+"."+constant.STICKY_KEY, strconv.FormatBool(v.Sticky))
		if len(v.RequestTimeout) != 0 {
			urlMap.Set("methods."+v.Name+"."+constant.TIMEOUT_KEY, v.RequestTimeout)
		}
	}

	return urlMap
}

// postProcessConfig asks registered ConfigPostProcessor to post-process the current ReferenceConfig.
func (c *ReferenceConfig) postProcessConfig(url *common.URL) {
	for _, p := range extension.GetConfigPostProcessors() {
		p.PostProcessReferenceConfig(url)
	}
}

// Refer ...
func (c *ReferenceConfig) Refer(_ interface{}) {
	cfgURL := common.NewURLWithOptions(
		common.WithPath(c.InterfaceName),
		common.WithProtocol(c.Protocol),
		common.WithParams(c.getUrlMap()),
		common.WithParamsValue(constant.BEAN_NAME_KEY, c.id),
	)
	if c.ForceTag {
		cfgURL.AddParam(constant.ForceUseTag, "true")
	}
	c.postProcessConfig(cfgURL)
	if c.Url != "" {
		// 1. user specified URL, could be peer-to-peer address, or register center's address.
		urlStrings := gxstrings.RegSplit(c.Url, "\\s*[;]+\\s*")
		for _, urlStr := range urlStrings {
			serviceUrl, err := common.NewURL(urlStr)
			if err != nil {
				panic(fmt.Sprintf("user specified URL %v refer error, error message is %v ", urlStr, err.Error()))
			}
			if serviceUrl.Protocol == constant.REGISTRY_PROTOCOL {
				serviceUrl.SubURL = cfgURL
				c.urls = append(c.urls, serviceUrl)
			} else {
				if serviceUrl.Path == "" {
					serviceUrl.Path = "/" + c.InterfaceName
				}
				// merge url need to do
				newUrl := common.MergeUrl(serviceUrl, cfgURL)
				c.urls = append(c.urls, newUrl)
			}
		}
	} else {
		// 2. assemble SubURL from register center's configuration mode
		c.urls = loadRegistries(c.Registry, consumerConfig.Registries, common.CONSUMER)

		// set url to regUrls
		for _, regUrl := range c.urls {
			regUrl.SubURL = cfgURL
		}
	}

	if len(c.urls) == 1 {
		c.invoker = extension.GetProtocol(c.urls[0].Protocol).Refer(c.urls[0])
	} else {
		invokers := make([]protocol.Invoker, 0, len(c.urls))
		var regUrl *common.URL
		for _, u := range c.urls {
			invokers = append(invokers, extension.GetProtocol(u.Protocol).Refer(u))
			if u.Protocol == constant.REGISTRY_PROTOCOL {
				regUrl = u
			}
		}

		// TODO(decouple from directory, config should not depend on directory module)
		var hitClu string
		if regUrl != nil {
			// for multi-subscription scenario, use 'zone-aware' policy by default
			hitClu = constant.ZONEAWARE_CLUSTER_NAME
		} else {
			// not a registry url, must be direct invoke.
			hitClu = constant.FAILOVER_CLUSTER_NAME
			if len(invokers) > 0 {
				u := invokers[0].GetUrl()
				if nil != &u {
					hitClu = u.GetParam(constant.CLUSTER_KEY, constant.ZONEAWARE_CLUSTER_NAME)
				}
			}
		}

		cluster := extension.GetCluster(hitClu)
		// If 'zone-aware' policy select, the invoker wrap sequence would be:
		// ZoneAwareClusterInvoker(StaticDirectory) ->
		// FailoverClusterInvoker(RegistryDirectory, routing happens here) -> Invoker
		c.invoker = cluster.Join(directory.NewStaticDirectory(invokers))
	}
	// publish consumer metadata
	//publishConsumerDefinition(cfgURL)
	// create proxy
	if c.Async {
		callback := GetCallback(c.id)
		c.pxy = extension.GetProxyFactory(consumerConfig.ProxyFactory).GetAsyncProxy(c.invoker, callback, cfgURL)
	} else {
		c.pxy = extension.GetProxyFactory(consumerConfig.ProxyFactory).GetProxy(c.invoker, cfgURL)
	}
}

// Implement
// @v is service provider implemented RPCService
func (c *ReferenceConfig) Implement(v common.RPCService) {
	c.pxy.Implement(v)
}
