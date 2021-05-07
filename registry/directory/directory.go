package directory

import (
	"fmt"
	"github.com/laz/dubbo-go/cluster"
	"github.com/laz/dubbo-go/cluster/directory"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/protocol"
	"github.com/laz/dubbo-go/registry"
	"net/url"
	"os"
	"sync"
)
import (
	perrors "github.com/pkg/errors"
)

func init() {
	extension.SetDefaultRegistryDirectory(NewRegistryDirectory)
}

type RegistryDirectory struct {
	directory.BaseDirectory
	cacheInvokers    []protocol.Invoker
	listenerLock     sync.Mutex
	serviceType      string
	registry         registry.Registry
	cacheInvokersMap *sync.Map // use sync.map
	consumerURL      *common.URL
	cacheOriginUrl   *common.URL
	/*configurators                  []config_center.Configurator
	consumerConfigurationListener  *consumerConfigurationListener
	referenceConfigurationListener *referenceConfigurationListener*/
	//serviceKey                     string
	//forbidden                      atomic.Bool
	registerLock sync.Mutex // this lock if for register
}

// NewRegistryDirectory will create a new RegistryDirectory
func NewRegistryDirectory(url *common.URL, registry registry.Registry) (cluster.Directory, error) {
	if url.SubURL == nil {
		return nil, perrors.Errorf("url is invalid, suburl can not be nil")
	}
	logger.Debugf("new RegistryDirectory for service :%s.", url.Key())
	dir := &RegistryDirectory{
		BaseDirectory:    directory.NewBaseDirectory(url),
		cacheInvokers:    []protocol.Invoker{},
		cacheInvokersMap: &sync.Map{},
		serviceType:      url.SubURL.Service(),
		registry:         registry,
	}

	dir.consumerURL = dir.getConsumerUrl(url.SubURL)

	go dir.subscribe(url.SubURL)
	return dir, nil
}

// List selected protocol invokers from the directory
func (dir *RegistryDirectory) List(invocation protocol.Invocation) []protocol.Invoker {
	invokers := dir.cacheInvokers

	return invokers

}

// Destroy method
func (dir *RegistryDirectory) Destroy() {
	// TODO:unregister & unsubscribe
	dir.BaseDirectory.Destroy(func() {
		invokers := dir.cacheInvokers
		dir.cacheInvokers = []protocol.Invoker{}
		for _, ivk := range invokers {
			ivk.Destroy()
		}
	})
}

// subscribe from registry
func (dir *RegistryDirectory) subscribe(url *common.URL) {

}

func (dir *RegistryDirectory) getConsumerUrl(c *common.URL) *common.URL {
	processID := fmt.Sprintf("%d", os.Getpid())
	localIP := common.GetLocalIp()

	params := url.Values{}
	c.RangeParams(func(key, value string) bool {
		params.Add(key, value)
		return true
	})

	params.Add("pid", processID)
	params.Add("ip", localIP)
	params.Add("protocol", c.Protocol)

	return common.NewURLWithOptions(common.WithProtocol("consumer"), common.WithIp(localIP), common.WithPath(c.Path),
		common.WithParams(params))
}
