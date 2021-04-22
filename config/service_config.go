package config

import "github.com/laz/dubbo-go/common"

// ServiceConfig is the configuration of the service provider
type ServiceConfig struct {
	Methods       []*MethodConfig `yaml:"methods"  json:"methods,omitempty" property:"methods"`
	InterfaceName string          `required:"true"  yaml:"interface"  json:"interface,omitempty" property:"interface"`
	id            string
	rpcService    common.RPCService
	Protocols     map[string]*ProtocolConfig
}

// Implement only store the @s and return
func (c *ServiceConfig) Implement(s common.RPCService) {
	c.rpcService = s
}

// Export exports the service
func (c *ServiceConfig) Export() error {

	return nil
}
