package config_center

import (
	"fmt"
	"github.com/laz/dubbo-go/remoting"
)

// ConfigurationListener for changing listener's event
type ConfigurationListener interface {
	// Process the notification event once there's any change happens on the config
	Process(*ConfigChangeEvent)
}

// ConfigChangeEvent for changing listener's event
type ConfigChangeEvent struct {
	Key        string
	Value      interface{}
	ConfigType remoting.EventType
}

func (c ConfigChangeEvent) String() string {
	return fmt.Sprintf("ConfigChangeEvent{key = %v , value = %v , changeType = %v}", c.Key, c.Value, c.ConfigType)
}
