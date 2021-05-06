package config

import (
	"bytes"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/common/yaml"
	"time"
)
import (
	perrors "github.com/pkg/errors"
)

const (
	MaxWheelTimeSpan = 900e9 // 900s, 15 minute
)

// ConsumerConfig is Consumer default configuration
type ConsumerConfig struct {
	BaseConfig   `yaml:",inline"`
	configCenter `yaml:"-"`
	Filter       string `yaml:"filter" json:"filter,omitempty" property:"filter"`
	// client
	Connect_Timeout string `default:"100ms"  yaml:"connect_timeout" json:"connect_timeout,omitempty" property:"connect_timeout"`
	ConnectTimeout  time.Duration

	Registry   *RegistryConfig            `yaml:"registry" json:"registry,omitempty" property:"registry"`
	Registries map[string]*RegistryConfig `default:"{}" yaml:"registries" json:"registries" property:"registries"`

	Request_Timeout string `yaml:"request_timeout" default:"5s" json:"request_timeout,omitempty" property:"request_timeout"`
	RequestTimeout  time.Duration
	ProxyFactory    string `yaml:"proxy_factory" default:"default" json:"proxy_factory,omitempty" property:"proxy_factory"`
	Check           *bool  `yaml:"check"  json:"check,omitempty" property:"check"`

	References   map[string]*ReferenceConfig `yaml:"references" json:"references,omitempty" property:"references"`
	ProtocolConf interface{}                 `yaml:"protocol_conf" json:"protocol_conf,omitempty" property:"protocol_conf"`
	FilterConf   interface{}                 `yaml:"filter_conf" json:"filter_conf,omitempty" property:"filter_conf"`

	ConfigType map[string]string `yaml:"config_type" json:"config_type,omitempty" property:"config_type"`
}

// ConsumerInit loads config file to init consumer config
func ConsumerInit(confConFile string) error {
	if confConFile == "" {
		return perrors.Errorf("application configure(consumer) file name is nil")
	}
	consumerConfig = &ConsumerConfig{}
	fileStream, err := yaml.UnmarshalYMLConfig(confConFile, consumerConfig)
	if err != nil {
		return perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	consumerConfig.fileStream = bytes.NewBuffer(fileStream)
	//set method interfaceId & interfaceName
	for k, v := range consumerConfig.References {
		//set id for reference
		for _, n := range consumerConfig.References[k].Methods {
			n.InterfaceName = v.InterfaceName
			n.InterfaceId = k
		}
	}
	if consumerConfig.Request_Timeout != "" {
		if consumerConfig.RequestTimeout, err = time.ParseDuration(consumerConfig.Request_Timeout); err != nil {
			return perrors.WithMessagef(err, "time.ParseDuration(Request_Timeout{%#v})", consumerConfig.Request_Timeout)
		}
		if consumerConfig.RequestTimeout >= time.Duration(MaxWheelTimeSpan) {
			return perrors.WithMessagef(err, "request_timeout %s should be less than %s",
				consumerConfig.Request_Timeout, time.Duration(MaxWheelTimeSpan))
		}
	}
	if consumerConfig.Connect_Timeout != "" {
		if consumerConfig.ConnectTimeout, err = time.ParseDuration(consumerConfig.Connect_Timeout); err != nil {
			return perrors.WithMessagef(err, "time.ParseDuration(Connect_Timeout{%#v})", consumerConfig.Connect_Timeout)
		}
	}

	logger.Debugf("consumer config{%#v}\n", consumerConfig)

	return nil
}

func configCenterRefreshConsumer() error {

	return nil
}
