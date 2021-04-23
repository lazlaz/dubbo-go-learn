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
