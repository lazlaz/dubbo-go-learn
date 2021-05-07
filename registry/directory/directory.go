package directory

import (
	"fmt"
	"github.com/laz/dubbo-go/cluster"
	"github.com/laz/dubbo-go/cluster/directory"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/protocol"
	"github.com/laz/dubbo-go/protocol/protocolwrapper"
	"github.com/laz/dubbo-go/registry"
	"github.com/laz/dubbo-go/remoting"
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
	//	configurators                  []config_center.Configurator
	consumerConfigurationListener  *consumerConfigurationListener
	referenceConfigurationListener *referenceConfigurationListener
	//serviceKey                     string
	//forbidden                      atomic.Bool
	registerLock sync.Mutex // this lock if for register
}

type referenceConfigurationListener struct {
	registry.BaseConfigurationListener
	directory *RegistryDirectory
	url       *common.URL
}
type consumerConfigurationListener struct {
	registry.BaseConfigurationListener
	listeners []registry.NotifyListener
	directory *RegistryDirectory
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
	dir.consumerConfigurationListener = newConsumerConfigurationListener(dir)
	go dir.subscribe(url.SubURL)
	return dir, nil
}
func newConsumerConfigurationListener(dir *RegistryDirectory) *consumerConfigurationListener {
	listener := &consumerConfigurationListener{directory: dir}

	return listener
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
	logger.Debugf("subscribe service :%s for RegistryDirectory.", url.Key())
	dir.consumerConfigurationListener.addNotifyListener(dir)

	if err := dir.registry.Subscribe(url, dir); err != nil {
		logger.Error("registry.Subscribe(url:%v, dir:%v) = error:%v", url, dir, err)
	}
}

// Notify monitor changes from registry,and update the cacheServices
func (dir *RegistryDirectory) Notify(event *registry.ServiceEvent) {
	if event == nil {
		return
	}
	dir.refreshInvokers(event)
}

// setNewInvokers groups the invokers from the cache first, then set the result to both directory and router chain.
func (dir *RegistryDirectory) setNewInvokers() {
	newInvokers := dir.toGroupInvokers()
	dir.listenerLock.Lock()
	defer dir.listenerLock.Unlock()
	dir.cacheInvokers = newInvokers
}
func (dir *RegistryDirectory) toGroupInvokers() []protocol.Invoker {
	var (
		_               error
		newInvokersList []protocol.Invoker
	)
	groupInvokersMap := make(map[string][]protocol.Invoker)

	dir.cacheInvokersMap.Range(func(key, value interface{}) bool {
		newInvokersList = append(newInvokersList, value.(protocol.Invoker))
		return true
	})

	for _, invoker := range newInvokersList {
		group := invoker.GetUrl().GetParam(constant.GROUP_KEY, "")

		groupInvokersMap[group] = append(groupInvokersMap[group], invoker)
	}
	groupInvokersList := make([]protocol.Invoker, 0, len(groupInvokersMap))

	// len is 1 it means no group setting ,so do not need cluster again
	for _, invokers := range groupInvokersMap {
		groupInvokersList = invokers
	}

	return groupInvokersList
}

// refreshInvokers refreshes service's events.
func (dir *RegistryDirectory) refreshInvokers(event *registry.ServiceEvent) {
	if event != nil {
		logger.Debugf("refresh invokers with %+v", event)
	} else {
		logger.Debug("refresh invokers with nil")
	}

	var oldInvoker protocol.Invoker
	if event != nil {
		oldInvoker, _ = dir.cacheInvokerByEvent(event)
	}
	dir.setNewInvokers()
	if oldInvoker != nil {
		oldInvoker.Destroy()
	}
}

// convertUrl processes override:// and router://
func (dir *RegistryDirectory) convertUrl(res *registry.ServiceEvent) *common.URL {
	ret := res.Service
	if ret.Protocol == constant.OVERRIDE_PROTOCOL || // 1.for override url in 2.6.x
		ret.GetParam(constant.CATEGORY_KEY, constant.DEFAULT_CATEGORY) == constant.CONFIGURATORS_CATEGORY {

		ret = nil
	} else if ret.Protocol == constant.ROUTER_PROTOCOL || // 2.for router
		ret.GetParam(constant.CATEGORY_KEY, constant.DEFAULT_CATEGORY) == constant.ROUTER_CATEGORY {
		ret = nil
	}
	return ret
}

func (dir *RegistryDirectory) overrideUrl(targetUrl *common.URL) {

}

// cacheInvoker will return abandoned Invoker,if no Invoker to be abandoned,return nil
func (dir *RegistryDirectory) cacheInvoker(url *common.URL) protocol.Invoker {
	dir.overrideUrl(dir.GetDirectoryUrl())
	referenceUrl := dir.GetDirectoryUrl().SubURL

	if url == nil && dir.cacheOriginUrl != nil {
		url = dir.cacheOriginUrl
	} else {
		dir.cacheOriginUrl = url
	}
	if url == nil {
		logger.Error("URL is nil ,pls check if service url is subscribe successfully!")
		return nil
	}
	// check the url's protocol is equal to the protocol which is configured in reference config or referenceUrl is not care about protocol
	if url.Protocol == referenceUrl.Protocol || referenceUrl.Protocol == "" {
		newUrl := common.MergeUrl(url, referenceUrl)
		dir.overrideUrl(newUrl)
		if v, ok := dir.doCacheInvoker(newUrl); ok {
			return v
		}
	}
	return nil
}

func (dir *RegistryDirectory) doCacheInvoker(newUrl *common.URL) (protocol.Invoker, bool) {
	key := newUrl.Key()
	if cacheInvoker, ok := dir.cacheInvokersMap.Load(key); !ok {
		logger.Debugf("service will be added in cache invokers: invokers url is  %s!", newUrl)
		newInvoker := extension.GetProtocol(protocolwrapper.FILTER).Refer(newUrl)
		if newInvoker != nil {
			dir.cacheInvokersMap.Store(key, newInvoker)
		} else {
			logger.Warnf("service will be added in cache invokers fail, result is null, invokers url is %+v", newUrl.String())
		}
	} else {
		// if cached invoker has the same URL with the new URL, then no need to re-refer, and no need to destroy
		// the old invoker.
		if common.GetCompareURLEqualFunc()(newUrl, cacheInvoker.(protocol.Invoker).GetUrl()) {
			return nil, true
		}

		logger.Debugf("service will be updated in cache invokers: new invoker url is %s, old invoker url is %s", newUrl, cacheInvoker.(protocol.Invoker).GetUrl())
		newInvoker := extension.GetProtocol(protocolwrapper.FILTER).Refer(newUrl)
		if newInvoker != nil {
			dir.cacheInvokersMap.Store(key, newInvoker)
			return cacheInvoker.(protocol.Invoker), true
		} else {
			logger.Warnf("service will be updated in cache invokers fail, result is null, invokers url is %+v", newUrl.String())
		}
	}
	return nil, false
}

// cacheInvokerByEvent caches invokers from the service event
func (dir *RegistryDirectory) cacheInvokerByEvent(event *registry.ServiceEvent) (protocol.Invoker, error) {
	// judge is override or others
	if event != nil {
		u := dir.convertUrl(event)
		switch event.Action {
		case remoting.EventTypeAdd, remoting.EventTypeUpdate:
			logger.Infof("selector add service url{%s}", event.Service)
			if u != nil && constant.ROUTER_PROTOCOL == u.Protocol {
				//dir.configRouters()
			}
			return dir.cacheInvoker(u), nil
		case remoting.EventTypeDel:
			logger.Infof("selector delete service url{%s}", event.Service)
			return dir.uncacheInvoker(u), nil
		default:
			return nil, fmt.Errorf("illegal event type: %v", event.Action)
		}
	}
	return nil, nil
}

// uncacheInvoker will return abandoned Invoker, if no Invoker to be abandoned, return nil
func (dir *RegistryDirectory) uncacheInvoker(url *common.URL) protocol.Invoker {
	return dir.uncacheInvokerWithKey(url.Key())
}

func (dir *RegistryDirectory) uncacheInvokerWithKey(key string) protocol.Invoker {
	logger.Debugf("service will be deleted in cache invokers: invokers key is  %s!", key)
	protocol.RemoveUrlKeyUnhealthyStatus(key)
	if cacheInvoker, ok := dir.cacheInvokersMap.Load(key); ok {
		dir.cacheInvokersMap.Delete(key)
		return cacheInvoker.(protocol.Invoker)
	}
	return nil
}

// NotifyAll notify the events that are complete Service Event List.
// After notify the address, the callback func will be invoked.
func (dir *RegistryDirectory) NotifyAll(events []*registry.ServiceEvent, callback func()) {
	go dir.refreshAllInvokers(events, callback)
}

// refreshAllInvokers the argument is the complete list of the service events,  we can safely assume any cached invoker
// not in the incoming list can be removed.  The Action of serviceEvent should be EventTypeUpdate.
func (dir *RegistryDirectory) refreshAllInvokers(events []*registry.ServiceEvent, callback func()) {
	var (
		oldInvokers []protocol.Invoker
		addEvents   []*registry.ServiceEvent
	)
	dir.overrideUrl(dir.GetDirectoryUrl())
	referenceUrl := dir.GetDirectoryUrl().SubURL

	// loop the events to check the Action should be EventTypeUpdate.
	for _, event := range events {
		if event.Action != remoting.EventTypeUpdate {
			panic("Your implements of register center is wrong, " +
				"please check the Action of ServiceEvent should be EventTypeUpdate")
		}
		// Originally it will Merge URL many times, now we just execute once.
		// MergeUrl is executed once and put the result into Event. After this, the key will get from Event.Key().
		newUrl := dir.convertUrl(event)
		newUrl = common.MergeUrl(newUrl, referenceUrl)
		dir.overrideUrl(newUrl)
		event.Update(newUrl)
	}
	// After notify all addresses, do some callback.
	defer callback()
	func() {
		// this lock is work at batch update of InvokeCache
		dir.registerLock.Lock()
		defer dir.registerLock.Unlock()
		// get need clear invokers from original invoker list
		dir.cacheInvokersMap.Range(func(k, v interface{}) bool {
			if !dir.eventMatched(k.(string), events) {
				// delete unused invoker from cache
				if invoker := dir.uncacheInvokerWithKey(k.(string)); invoker != nil {
					oldInvokers = append(oldInvokers, invoker)
				}
			}
			return true
		})
		// get need add invokers from events
		for _, event := range events {
			// Get the key from Event.Key()
			if _, ok := dir.cacheInvokersMap.Load(event.Key()); !ok {
				addEvents = append(addEvents, event)
			}
		}
		// loop the updateEvents
		for _, event := range addEvents {
			logger.Debugf("registry update, result{%s}", event)
			if event.Service != nil {
				logger.Infof("selector add service url{%s}", event.Service.String())
			}
			if event != nil && event.Service != nil && constant.ROUTER_PROTOCOL == event.Service.Protocol {
				//	dir.configRouters()
			}
			if oldInvoker, _ := dir.doCacheInvoker(event.Service); oldInvoker != nil {
				oldInvokers = append(oldInvokers, oldInvoker)
			}
		}
	}()
	dir.setNewInvokers()
	// destroy unused invokers
	for _, invoker := range oldInvokers {
		go invoker.Destroy()
	}
}

// eventMatched checks if a cached invoker appears in the incoming invoker list, if no, then it is safe to remove.
func (dir *RegistryDirectory) eventMatched(key string, events []*registry.ServiceEvent) bool {
	for _, event := range events {
		if dir.invokerCacheKey(event) == key {
			return true
		}
	}
	return false
}

// invokerCacheKey generates the key in the cache for a given ServiceEvent.
func (dir *RegistryDirectory) invokerCacheKey(event *registry.ServiceEvent) string {
	// If the url is merged, then return Event.Key() directly.
	if event.Updated() {
		return event.Key()
	}
	referenceUrl := dir.GetDirectoryUrl().SubURL
	newUrl := common.MergeUrl(event.Service, referenceUrl)
	event.Update(newUrl)
	return newUrl.Key()
}

func (l *consumerConfigurationListener) addNotifyListener(listener registry.NotifyListener) {
	l.listeners = append(l.listeners, listener)
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
