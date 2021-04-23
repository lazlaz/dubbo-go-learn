package extension

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/registry"
)

var (
	registrys = make(map[string]func(config *common.URL) (registry.Registry, error))
)

// SetRegistry sets the registry extension with @name
func SetRegistry(name string, v func(_ *common.URL) (registry.Registry, error)) {
	registrys[name] = v
}

// GetRegistry finds the registry extension with @name
func GetRegistry(name string, config *common.URL) (registry.Registry, error) {
	if registrys[name] == nil {
		panic("registry for " + name + " does not exist. please make sure that you have imported the package `github.com/apache/dubbo-go/registry/" + name + "`.")
	}
	return registrys[name](config)

}
