package config

// ProtocolConfig is protocol configuration
type ProtocolConfig struct {
	Name string `required:"true" yaml:"name"  json:"name,omitempty" property:"name"`
	Ip   string `required:"true" yaml:"ip"  json:"ip,omitempty" property:"ip"`
	Port string `required:"true" yaml:"port"  json:"port,omitempty" property:"port"`
}
