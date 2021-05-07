package dubbo

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/config"
	"github.com/laz/dubbo-go/protocol"
	"github.com/laz/dubbo-go/remoting"
	"sync"
	"time"
)

// DubboInvoker is implement of protocol.Invoker. A dubboInvoker refers to one service and ip.
type DubboInvoker struct {
	protocol.BaseInvoker
	// the exchange layer, it is focus on network communication.
	clientGuard *sync.RWMutex
	client      *remoting.ExchangeClient
	quitOnce    sync.Once
	// timeout for service(interface) level.
	timeout time.Duration
}

// NewDubboInvoker constructor
func NewDubboInvoker(url *common.URL, client *remoting.ExchangeClient) *DubboInvoker {
	requestTimeout := config.GetConsumerConfig().RequestTimeout

	requestTimeoutStr := url.GetParam(constant.TIMEOUT_KEY, config.GetConsumerConfig().Request_Timeout)
	if t, err := time.ParseDuration(requestTimeoutStr); err == nil {
		requestTimeout = t
	}
	di := &DubboInvoker{
		BaseInvoker: *protocol.NewBaseInvoker(url),
		clientGuard: &sync.RWMutex{},
		client:      client,
		timeout:     requestTimeout,
	}

	return di
}
