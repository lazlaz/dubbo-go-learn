package config

import (
	"flag"
	"fmt"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/logger"
	"log"
	"os"
	"time"
)

var (
	consumerConfig *ConsumerConfig
	providerConfig *ProviderConfig
	baseConfig     *BaseConfig
	confRouterFile string
	sslEnabled     = false
	maxWait        = 3
)

func init() {
	var (
		confConFile string
		confProFile string
	)
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&confConFile, "conConf", os.Getenv(constant.CONF_CONSUMER_FILE_PATH), "default client config path")
	fs.StringVar(&confProFile, "proConf", os.Getenv(constant.CONF_PROVIDER_FILE_PATH), "default server config path")
	fs.StringVar(&confRouterFile, "rouConf", os.Getenv(constant.CONF_ROUTER_FILE_PATH), "default router config path")
	fs.Parse(os.Args[1:])
	for len(fs.Args()) != 0 {
		fs.Parse(fs.Args()[1:])
	}
	if errCon := ConsumerInit(confConFile); errCon != nil {
		log.Printf("[consumerInit] %#v", errCon)
		consumerConfig = nil
	} else {
		// Even though baseConfig has been initialized, we override it
		// because we think read from config file is correct config
		baseConfig = &consumerConfig.BaseConfig
	}
	if errPro := ProviderInit(confProFile); errPro != nil {
		log.Printf("[providerInit] %#v", errPro)
		providerConfig = nil
	} else {
		// Even though baseConfig has been initialized, we override it
		// because we think read from config file is correct config
		baseConfig = &providerConfig.BaseConfig
	}
}

func Load() {
	//初始化引用配置
	loadConsumerConfig()
	//初始化服务配置
	loadProviderConfig()
}

func loadConsumerConfig() {
	if consumerConfig == nil {
		logger.Warnf("consumerConfig is nil!")
		return
	}

	checkApplicationName(consumerConfig.ApplicationConfig)
	if err := configCenterRefreshConsumer(); err != nil {
		logger.Errorf("[consumer config center refresh] %#v", err)
	}
	checkRegistries(consumerConfig.Registries, consumerConfig.Registry)
	for key, ref := range consumerConfig.References {

		rpcService := GetConsumerService(key)
		if rpcService == nil {
			logger.Warnf("%s does not exist!", key)
			continue
		}
		ref.id = key
		ref.Refer(rpcService)
		ref.Implement(rpcService)
	}

	// wait for invoker is available, if wait over default 3s, then panic
	var count int
	for {
		checkok := true
		for _, refconfig := range consumerConfig.References {
			if (refconfig.Check != nil && *refconfig.Check) ||
				(refconfig.Check == nil && consumerConfig.Check != nil && *consumerConfig.Check) ||
				(refconfig.Check == nil && consumerConfig.Check == nil) { // default to true

				if refconfig.invoker != nil && !refconfig.invoker.IsAvailable() {
					checkok = false
					count++
					if count > maxWait {
						errMsg := fmt.Sprintf("Failed to check the status of the service %v. No provider available for the service to the consumer use dubbo version %v", refconfig.InterfaceName, constant.Version)
						logger.Error(errMsg)
						panic(errMsg)
					}
					time.Sleep(time.Second * 1)
					break
				}
				if refconfig.invoker == nil {
					logger.Warnf("The interface %s invoker not exist, may you should check your interface config.", refconfig.InterfaceName)
				}
			}
		}
		if checkok {
			break
		}
	}
}

// GetConsumerConfig find the consumer config
// if not found, create new one
// we use double-check to reduce race condition
// In general, it will be locked 0 or 1 time.
// So you don't need to worry about the race condition
func GetConsumerConfig() ConsumerConfig {
	if consumerConfig == nil {
		if consumerConfig == nil {
			return ConsumerConfig{}
		}
	}
	return *consumerConfig
}
func loadProviderConfig() {
	if providerConfig == nil {
		logger.Warnf("providerConfig is nil!")
		return
	}

	checkApplicationName(providerConfig.ApplicationConfig)
	if err := configCenterRefreshProvider(); err != nil {
		logger.Errorf("[provider config center refresh] %#v", err)
	}
	checkRegistries(providerConfig.Registries, providerConfig.Registry)

	for key, svs := range providerConfig.Services {
		rpcService := GetProviderService(key)
		if rpcService == nil {
			logger.Warnf("%s does not exist!", key)
			continue
		}
		svs.id = key
		svs.Implement(rpcService)
		svs.Protocols = providerConfig.Protocols
		if err := svs.Export(); err != nil {
			panic(fmt.Sprintf("service %s export failed! err: %#v", key, err))
		}
	}
	//注册服务实例
	//registerServiceInstance()
}
func checkRegistries(registries map[string]*RegistryConfig, singleRegistry *RegistryConfig) {
	if len(registries) == 0 && singleRegistry != nil {
		registries[constant.DEFAULT_KEY] = singleRegistry
	}
}
func checkApplicationName(config *ApplicationConfig) {
	if config == nil || len(config.Name) == 0 {
		errMsg := "application config must not be nil, pls check your configuration"
		logger.Errorf(errMsg)
		panic(errMsg)
	}
}

// GetProviderConfig find the provider config
// if not found, create new one
func GetProviderConfig() ProviderConfig {
	if providerConfig == nil {
		if providerConfig == nil {
			return ProviderConfig{}
		}
	}
	return *providerConfig
}

func GetSslEnabled() bool {
	return sslEnabled
}
