package extension

import "github.com/laz/dubbo-go/protocol"

var (
	protocols = make(map[string]func() protocol.Protocol)
)

// GetProtocol finds the protocol extension with @name
func GetProtocol(name string) protocol.Protocol {
	if protocols[name] == nil {
		panic("protocol for " + name + " is not existing, make sure you have import the package.")
	}
	return protocols[name]()
}
