package config

type ApplicationConfig struct {
	Organization string `yaml:"organization" json:"organization,omitempty" property:"organization"`
	Name         string `yaml:"name" json:"name,omitempty" property:"name"`
	Module       string `yaml:"module" json:"module,omitempty" property:"module"`
	Version      string `yaml:"version" json:"version,omitempty" property:"version"`
	Owner        string `yaml:"owner" json:"owner,omitempty" property:"owner"`
	Environment  string `yaml:"environment" json:"environment,omitempty" property:"environment"`
	// the metadata type. remote or local
	MetadataType string `default:"local" yaml:"metadataType" json:"metadataType,omitempty" property:"metadataType"`
}
