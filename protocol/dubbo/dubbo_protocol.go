package dubbo

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/protocol"
	"github.com/laz/dubbo-go/remoting"
	"sync"
)

const (
	// DUBBO is dubbo protocol name
	DUBBO = "dubbo"
)

func init() {
	extension.SetProtocol(DUBBO, GetProtocol)
}

var (
	dubboProtocol *DubboProtocol
)

// It support dubbo protocol. It implements Protocol interface for dubbo protocol.
type DubboProtocol struct {
	protocol.BaseProtocol
	// It is store relationship about serviceKey(group/interface:version) and ExchangeServer
	// The ExchangeServer is introduced to replace of Server. Because Server is depend on getty directly.
	serverMap  map[string]*remoting.ExchangeServer
	serverLock sync.Mutex
}

// GetProtocol get a single dubbo protocol.
func GetProtocol() protocol.Protocol {
	if dubboProtocol == nil {
		dubboProtocol = NewDubboProtocol()
	}
	return dubboProtocol
}

// Export export dubbo service.
func (dp *DubboProtocol) Export(invoker protocol.Invoker) protocol.Exporter {
	return nil

}

// Refer create dubbo service reference.
func (dp *DubboProtocol) Refer(url *common.URL) protocol.Invoker {
	return nil
}

// Destroy destroy dubbo service.
func (dp *DubboProtocol) Destroy() {
	logger.Infof("DubboProtocol destroy.")

}

// NewDubboProtocol create a dubbo protocol.
func NewDubboProtocol() *DubboProtocol {
	return &DubboProtocol{
		BaseProtocol: protocol.NewBaseProtocol(),
		serverMap:    make(map[string]*remoting.ExchangeServer),
	}
}
