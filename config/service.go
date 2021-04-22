package config

import "github.com/laz/dubbo-go/common"

var (
	conServices = map[string]common.RPCService{} // service name -> service
	proServices = map[string]common.RPCService{} // service name -> service
)

// GetProviderService gets ProviderService by @name
func GetProviderService(name string) common.RPCService {
	return proServices[name]
}
