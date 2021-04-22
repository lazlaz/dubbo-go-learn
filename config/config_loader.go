package config

import (
	"flag"
	"fmt"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/logger"
	"log"
	"os"
)

var (
	providerConfig *ProviderConfig
	baseConfig     *BaseConfig
	confRouterFile string
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
	//初始化服务配置
	loadProviderConfig()
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
