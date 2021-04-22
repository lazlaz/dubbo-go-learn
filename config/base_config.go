package config

import "bytes"

type BaseConfig struct {
	fileStream *bytes.Buffer
	// application config
	ApplicationConfig *ApplicationConfig `yaml:"application" json:"application,omitempty" property:"application"`
}
