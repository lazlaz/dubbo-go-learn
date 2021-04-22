package config

import "strings"

// ProtocolConfig is protocol configuration
type ProtocolConfig struct {
	Name string `required:"true" yaml:"name"  json:"name,omitempty" property:"name"`
	Ip   string `required:"true" yaml:"ip"  json:"ip,omitempty" property:"ip"`
	Port string `required:"true" yaml:"port"  json:"port,omitempty" property:"port"`
}

func loadProtocol(protocolsIds string, protocols map[string]*ProtocolConfig) []*ProtocolConfig {
	returnProtocols := make([]*ProtocolConfig, 0, len(protocols))
	for _, v := range strings.Split(protocolsIds, ",") {
		for k, protocol := range protocols {
			if v == k {
				returnProtocols = append(returnProtocols, protocol)
			}
		}
	}
	return returnProtocols
}
