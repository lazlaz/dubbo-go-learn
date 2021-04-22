package config

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"net/url"
)
import (
	perrors "github.com/pkg/errors"
)

type ConfigCenterConfig struct {
	//context       context.Context
	Protocol      string `required:"true"  yaml:"protocol"  json:"protocol,omitempty"`
	Address       string `yaml:"address" json:"address,omitempty"`
	Cluster       string `yaml:"cluster" json:"cluster,omitempty"`
	Group         string `default:"dubbo" yaml:"group" json:"group,omitempty"`
	Username      string `yaml:"username" json:"username,omitempty"`
	Password      string `yaml:"password" json:"password,omitempty"`
	LogDir        string `yaml:"log_dir" json:"log_dir,omitempty"`
	ConfigFile    string `default:"dubbo.properties" yaml:"config_file"  json:"config_file,omitempty"`
	Namespace     string `default:"dubbo" yaml:"namespace"  json:"namespace,omitempty"`
	AppConfigFile string `default:"dubbo.properties" yaml:"app_config_file"  json:"app_config_file,omitempty"`
	AppId         string `default:"dubbo" yaml:"app_id"  json:"app_id,omitempty"`
	TimeoutStr    string `yaml:"timeout"  json:"timeout,omitempty"`
	RemoteRef     string `required:"false"  yaml:"remote_ref"  json:"remote_ref,omitempty"`
	//timeout       time.Duration
}
type configCenter struct {
}

//为configCenter结构体绑定方法
// toURL will compatible with baseConfig.ConfigCenterConfig.Address and baseConfig.ConfigCenterConfig.RemoteRef before 1.6.0
// After 1.6.0 will not compatible, only baseConfig.ConfigCenterConfig.RemoteRef
func (b *configCenter) toURL(baseConfig BaseConfig) (*common.URL, error) {
	if len(baseConfig.ConfigCenterConfig.Address) > 0 {
		return common.NewURL(baseConfig.ConfigCenterConfig.Address,
			common.WithProtocol(baseConfig.ConfigCenterConfig.Protocol), common.WithParams(baseConfig.ConfigCenterConfig.GetUrlMap()))
	}

	remoteRef := baseConfig.ConfigCenterConfig.RemoteRef
	rc, ok := baseConfig.GetRemoteConfig(remoteRef)

	if !ok {
		return nil, perrors.New("Could not find out the remote ref config, name: " + remoteRef)
	}

	newURL, err := rc.toURL()
	if err == nil {
		newURL.SetParams(baseConfig.ConfigCenterConfig.GetUrlMap())
	}
	return newURL, err
}

// GetUrlMap gets url map from ConfigCenterConfig
func (c *ConfigCenterConfig) GetUrlMap() url.Values {
	urlMap := url.Values{}
	urlMap.Set(constant.CONFIG_NAMESPACE_KEY, c.Namespace)
	urlMap.Set(constant.CONFIG_GROUP_KEY, c.Group)
	urlMap.Set(constant.CONFIG_CLUSTER_KEY, c.Cluster)
	urlMap.Set(constant.CONFIG_APP_ID_KEY, c.AppId)
	urlMap.Set(constant.CONFIG_LOG_DIR_KEY, c.LogDir)
	return urlMap
}

// startConfigCenter will start the config center.
// it will prepare the environment
func (b *configCenter) startConfigCenter(baseConfig BaseConfig) error {
	//newUrl, err := b.toURL(baseConfig)
	//if err != nil {
	//	return err
	//}
	//if err = b.prepareEnvironment(baseConfig, newUrl); err != nil {
	//	return perrors.WithMessagef(err, "start config center error!")
	//}
	// c.fresh()
	return nil
}
