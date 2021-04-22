package config

type ApplicationConfig struct {
	Name string `yaml:"name" json:"name,omitempty" property:"name"`
}
