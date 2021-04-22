package config

// ServiceConfig is the configuration of the service provider
type ServiceConfig struct {
	Methods       []*MethodConfig `yaml:"methods"  json:"methods,omitempty" property:"methods"`
	InterfaceName string          `required:"true"  yaml:"interface"  json:"interface,omitempty" property:"interface"`
}
