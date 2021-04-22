package config

import "bytes"

type BaseConfig struct {
	fileStream *bytes.Buffer
	// application config
	ApplicationConfig  *ApplicationConfig  `yaml:"application" json:"application,omitempty" property:"application"`
	ConfigCenterConfig *ConfigCenterConfig `yaml:"config_center" json:"config_center,omitempty"`
	//prefix              string
	fatherConfig interface{}
	// since 1.5.0 version
	Remotes map[string]*RemoteConfig `yaml:"remote" json:"remote,omitempty"`
}

// GetRemoteConfig will return the remote's config with the name if found
func (c *BaseConfig) GetRemoteConfig(name string) (config *RemoteConfig, ok bool) {
	config, ok = c.Remotes[name]
	return
}
