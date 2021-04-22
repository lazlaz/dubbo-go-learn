package config

import (
	"bytes"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/laz/dubbo-go/common/yaml"
)

type ProviderConfig struct {
	BaseConfig `yaml:",inline"`
	Services   map[string]*ServiceConfig `yaml:"services" json:"services,omitempty" property:"services"`
	ConfigType map[string]string         `yaml:"config_type" json:"config_type,omitempty" property:"config_type"`
}

func ProviderInit(confProFile string) error {
	if len(confProFile) == 0 {
		return perrors.Errorf("application configure(provider) file name is nil")
	}
	providerConfig = &ProviderConfig{}
	fileStream, err := yaml.UnmarshalYMLConfig(confProFile, providerConfig)
	if err != nil {
		return perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}

	providerConfig.fileStream = bytes.NewBuffer(fileStream)
	// set method interfaceId & interfaceName
	for k, v := range providerConfig.Services {
		// set id for reference
		for _, n := range providerConfig.Services[k].Methods {
			n.InterfaceName = v.InterfaceName
			n.InterfaceId = k
		}
	}

	return nil
}
