package nacos

import (
	"bytes"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"strconv"
	"strings"
	"time"
)
import (
	nacosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	perrors "github.com/pkg/errors"
)

var (
	localIP = ""
)

const (
	// RegistryConnDelay registry connection delay
	RegistryConnDelay = 3
)

type nacosRegistry struct {
	*common.URL
	namingClient naming_client.INamingClient
	registryUrls []*common.URL
}

func init() {
	localIP = common.GetLocalIp()
	extension.SetRegistry(constant.NACOS_KEY, newNacosRegistry)
}

// GetUrl gets its registration URL
func (nr *nacosRegistry) GetUrl() *common.URL {
	return nr.URL
}

// IsAvailable determines nacos registry center whether it is available
func (nr *nacosRegistry) IsAvailable() bool {
	// TODO
	return true
}

// nolint
func (nr *nacosRegistry) Destroy() {

}

// UnRegister
func (nr *nacosRegistry) UnRegister(conf *common.URL) error {
	return perrors.New("UnRegister is not support in nacosRegistry")
}

// subscribe from registry
func (nr *nacosRegistry) Subscribe(url *common.URL, notifyListener registry.NotifyListener) error {
	role, _ := strconv.Atoi(nr.URL.GetParam(constant.ROLE_KEY, ""))
	if role != common.CONSUMER {
		return nil
	}

	for {
		if !nr.IsAvailable() {
			logger.Warnf("event listener game over.")
			return perrors.New("nacosRegistry is not available.")
		}

		listener, err := nr.subscribe(url)
		if err != nil {
			if !nr.IsAvailable() {
				logger.Warnf("event listener game over.")
				return err
			}
			logger.Warnf("getListener() = err:%v", perrors.WithStack(err))
			time.Sleep(time.Duration(RegistryConnDelay) * time.Second)
			continue
		}

		for {
			serviceEvent, err := listener.Next()
			if err != nil {
				logger.Warnf("Selector.watch() = error{%v}", perrors.WithStack(err))
				listener.Close()
				return err
			}

			logger.Infof("update begin, service event: %v", serviceEvent.String())
			notifyListener.Notify(serviceEvent)
		}

	}
}
func (nr *nacosRegistry) subscribe(conf *common.URL) (registry.Listener, error) {
	return NewNacosListener(conf, nr.namingClient)
}

// UnSubscribe :
func (nr *nacosRegistry) UnSubscribe(url *common.URL, notifyListener registry.NotifyListener) error {
	return perrors.New("UnSubscribe not support in nacosRegistry")
}

// Register will register the service @url to its nacos registry center
func (nr *nacosRegistry) Register(url *common.URL) error {
	serviceName := getServiceName(url)
	//创建注册到nacos中的信息
	param := createRegisterParam(url, serviceName)
	isRegistry, err := nr.namingClient.RegisterInstance(param)
	if err != nil {
		return err
	}
	if !isRegistry {
		return perrors.New("registry [" + serviceName + "] to  nacos failed")
	}
	nr.registryUrls = append(nr.registryUrls, url)
	return nil
}

func createRegisterParam(url *common.URL, serviceName string) vo.RegisterInstanceParam {
	category := getCategory(url)
	params := make(map[string]string)

	url.RangeParams(func(key, value string) bool {
		params[key] = value
		return true
	})

	params[constant.NACOS_CATEGORY_KEY] = category
	params[constant.NACOS_PROTOCOL_KEY] = url.Protocol
	params[constant.NACOS_PATH_KEY] = url.Path
	if len(url.Ip) == 0 {
		url.Ip = localIP
	}
	if len(url.Port) == 0 || url.Port == "0" {
		url.Port = "80"
	}
	port, _ := strconv.Atoi(url.Port)
	instance := vo.RegisterInstanceParam{
		Ip:          url.Ip,
		Port:        uint64(port),
		Metadata:    params,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		ServiceName: serviceName,
	}
	return instance
}

func getCategory(url *common.URL) string {
	role, _ := strconv.Atoi(url.GetParam(constant.ROLE_KEY, strconv.Itoa(constant.NACOS_DEFAULT_ROLETYPE)))
	category := common.DubboNodes[role]
	return category
}

func getServiceName(url *common.URL) string {
	var buffer bytes.Buffer

	buffer.Write([]byte(getCategory(url)))
	appendParam(&buffer, url, constant.INTERFACE_KEY)
	appendParam(&buffer, url, constant.VERSION_KEY)
	appendParam(&buffer, url, constant.GROUP_KEY)
	return buffer.String()
}
func appendParam(target *bytes.Buffer, url *common.URL, key string) {
	value := url.GetParam(key, "")
	if strings.TrimSpace(value) != "" {
		target.Write([]byte(constant.NACOS_SERVICE_NAME_SEPARATOR))
		target.Write([]byte(value))
	}
}

// newNacosRegistry will create new instance
func newNacosRegistry(url *common.URL) (registry.Registry, error) {
	nacosConfig, err := getNacosConfig(url)
	if err != nil {
		return &nacosRegistry{}, err
	}
	client, err := clients.CreateNamingClient(nacosConfig)
	if err != nil {
		return &nacosRegistry{}, err
	}
	tmpRegistry := &nacosRegistry{
		URL:          url,
		namingClient: client,
		registryUrls: []*common.URL{},
	}
	return tmpRegistry, nil
}

// getNacosConfig will return the nacos config
// TODO support RemoteRef
func getNacosConfig(url *common.URL) (map[string]interface{}, error) {
	if url == nil {
		return nil, perrors.New("url is empty!")
	}
	if len(url.Location) == 0 {
		return nil, perrors.New("url.location is empty!")
	}
	configMap := make(map[string]interface{}, 2)

	addresses := strings.Split(url.Location, ",")
	serverConfigs := make([]nacosConstant.ServerConfig, 0, len(addresses))
	for _, addr := range addresses {
		ip, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, perrors.WithMessagef(err, "split [%s] ", addr)
		}
		port, _ := strconv.Atoi(portStr)
		serverConfigs = append(serverConfigs, nacosConstant.ServerConfig{
			IpAddr: ip,
			Port:   uint64(port),
		})
	}
	configMap[nacosConstant.KEY_SERVER_CONFIGS] = serverConfigs

	var clientConfig nacosConstant.ClientConfig
	timeout, err := time.ParseDuration(url.GetParam(constant.REGISTRY_TIMEOUT_KEY, constant.DEFAULT_REG_TIMEOUT))
	if err != nil {
		return nil, err
	}
	clientConfig.TimeoutMs = uint64(timeout.Seconds() * 1000)
	clientConfig.ListenInterval = 2 * clientConfig.TimeoutMs
	clientConfig.CacheDir = url.GetParam(constant.NACOS_CACHE_DIR_KEY, "")
	clientConfig.LogDir = url.GetParam(constant.NACOS_LOG_DIR_KEY, "")
	clientConfig.Endpoint = url.GetParam(constant.NACOS_ENDPOINT, "")
	clientConfig.NamespaceId = url.GetParam(constant.NACOS_NAMESPACE_ID, "")

	//enable local cache when nacos can not connect.
	notLoadCache, err := strconv.ParseBool(url.GetParam(constant.NACOS_NOT_LOAD_LOCAL_CACHE, "false"))
	if err != nil {
		logger.Errorf("ParseBool - error: %v", err)
		notLoadCache = false
	}
	clientConfig.NotLoadCacheAtStart = notLoadCache

	configMap[nacosConstant.KEY_CLIENT_CONFIG] = clientConfig

	return configMap, nil
}
