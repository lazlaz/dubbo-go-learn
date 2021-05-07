package extension

import (
	"github.com/laz/dubbo-go/cluster"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/registry"
)

type registryDirectory func(url *common.URL, registry registry.Registry) (cluster.Directory, error)

var defaultRegistry registryDirectory

// SetDefaultRegistryDirectory sets the default registryDirectory
func SetDefaultRegistryDirectory(v registryDirectory) {
	defaultRegistry = v
}

// GetDefaultRegistryDirectory finds the registryDirectory with url and registry
func GetDefaultRegistryDirectory(config *common.URL, registry registry.Registry) (cluster.Directory, error) {
	if defaultRegistry == nil {
		panic("registry directory is not existing, make sure you have import the package.")
	}
	return defaultRegistry(config, registry)
}
