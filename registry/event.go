package registry

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/remoting"
)

// ServiceEvent includes create, update, delete event
type ServiceEvent struct {
	Action  remoting.EventType
	Service *common.URL
	// store the key for Service.Key()
	key string
	// If the url is updated, such as Merged.
	updated bool
}

// Update() update the url with the merged URL. Work with Updated() can reduce the process of some merging URL.
func (e *ServiceEvent) Update(url *common.URL) {
	e.Service = url
	e.updated = true
}

// Updated() check if the url is updated.
// If the serviceEvent is updated, then it don't need merge url again.
func (e *ServiceEvent) Updated() bool {
	return e.updated
}

// Key() generate the key for service.Key(). It is cached once.
func (e *ServiceEvent) Key() string {
	if len(e.key) > 0 {
		return e.key
	}
	e.key = e.Service.Key()
	return e.key
}
