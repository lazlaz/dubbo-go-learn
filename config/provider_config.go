package config

import (
	"bytes"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/laz/dubbo-go/common/yaml"
)

type ProviderConfig struct {
	BaseConfig   `yaml:",inline"`
	Services     map[string]*ServiceConfig  `yaml:"services" json:"services,omitempty" property:"services"`
	ConfigType   map[string]string          `yaml:"config_type" json:"config_type,omitempty" property:"config_type"`
	Registries   map[string]*RegistryConfig `default:"{}" yaml:"registries" json:"registries" property:"registries"`
	Registry     *RegistryConfig            `yaml:"registry" json:"registry,omitempty" property:"registry"`
	Protocols    map[string]*ProtocolConfig `yaml:"protocols" json:"protocols,omitempty" property:"protocols"`
	Filter       string                     `yaml:"filter" json:"filter,omitempty" property:"filter"`
	ProxyFactory string                     `yaml:"proxy_factory" default:"default" json:"proxy_factory,omitempty" property:"proxy_factory"`
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

func configCenterRefreshProvider() error {
	// fresh it
	/*if providerConfig.ConfigCenterConfig != nil {
		providerConfig.fatherConfig = providerConfig
		if err := providerConfig.startConfigCenter((*providerConfig).BaseConfig); err != nil {
			return perrors.Errorf("start config center error , error message is {%v}", perrors.WithStack(err))
		}
		providerConfig.fresh()
	}*/
	return nil
}
