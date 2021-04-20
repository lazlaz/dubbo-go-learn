package config

import (
	"github.com/laz/dubbo-go/common/logger"
)

var (
	providerConfig *ProviderConfig
)

func Load() {
	//初始化服务配置
	loadProviderConfig()
}

func loadProviderConfig() {
	if providerConfig == nil {
		logger.Warnf("providerConfig is nil!")
		return
	}
}
