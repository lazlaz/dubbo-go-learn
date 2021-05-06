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

// SetProviderService is called by init() of implement of RPCService
func SetProviderService(service common.RPCService) {
	proServices[service.Reference()] = service
}

// GetConsumerService gets ConsumerService by @name
func GetConsumerService(name string) common.RPCService {
	return conServices[name]
}

// SetConsumerService is called by init() of implement of RPCService
func SetConsumerService(service common.RPCService) {
	conServices[service.Reference()] = service
}

// GetCallback gets CallbackResponse by @name
func GetCallback(name string) func(response common.CallbackResponse) {
	service := GetConsumerService(name)
	if sv, ok := service.(common.AsyncCallbackService); ok {
		return sv.CallBack
	}
	return nil
}
