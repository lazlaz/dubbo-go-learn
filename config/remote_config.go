package config

import "github.com/laz/dubbo-go/common"
import (
	perrors "github.com/pkg/errors"
)

type RemoteConfig struct {
	Protocol   string            `yaml:"protocol"  json:"protocol,omitempty"`
	Address    string            `yaml:"address" json:"address,omitempty"`
	TimeoutStr string            `default:"5s" yaml:"timeout" json:"timeout,omitempty"`
	Username   string            `yaml:"username" json:"username,omitempty" property:"username"`
	Password   string            `yaml:"password" json:"password,omitempty"  property:"password"`
	Params     map[string]string `yaml:"params" json:"params,omitempty"`
}

func (rc *RemoteConfig) toURL() (*common.URL, error) {
	if len(rc.Protocol) == 0 {
		return nil, perrors.Errorf("Must provide protocol in RemoteConfig.")
	}
	return common.NewURL(rc.Address,
		common.WithUsername(rc.Username),
		common.WithPassword(rc.Password),
		common.WithLocation(rc.Address),
		common.WithProtocol(rc.Protocol),
	)
}
