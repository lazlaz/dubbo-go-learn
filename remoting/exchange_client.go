package remoting

import (
	"errors"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/logger"
	"time"
)
import (
	uatomic "go.uber.org/atomic"
)

// It is interface of client for network communication.
// If you use getty as network communication, you should define GettyClient that implements this interface.
type Client interface {
	SetExchangeClient(client *ExchangeClient)
	// connect url
	Connect(url *common.URL) error
	// close
	Close()
	// send request to server.
	Request(request *Request, timeout time.Duration, response *PendingResponse) error
	// check if the client is still available
	IsAvailable() bool
}

// This is abstraction level. it is like facade.
type ExchangeClient struct {
	// connect server timeout
	ConnectTimeout time.Duration
	// to dial server address. The format: ip:port
	address string
	// the client that will deal with the transport. It is interface, and it will use gettyClient by default.
	client Client
	// the tag for init.
	init bool
	// the number of service using the exchangeClient
	activeNum uatomic.Uint32
}

// create ExchangeClient
func NewExchangeClient(url *common.URL, client Client, connectTimeout time.Duration, lazyInit bool) *ExchangeClient {
	exchangeClient := &ExchangeClient{
		ConnectTimeout: connectTimeout,
		address:        url.Location,
		client:         client,
	}
	client.SetExchangeClient(exchangeClient)
	if !lazyInit {
		if err := exchangeClient.doInit(url); err != nil {
			return nil
		}
	}
	exchangeClient.IncreaseActiveNumber()
	return exchangeClient
}

// increase number of service using client
func (client *ExchangeClient) IncreaseActiveNumber() uint32 {
	return client.activeNum.Add(1)
}

func (cl *ExchangeClient) doInit(url *common.URL) error {
	if cl.init {
		return nil
	}
	if cl.client.Connect(url) != nil {
		//retry for a while
		time.Sleep(100 * time.Millisecond)
		if cl.client.Connect(url) != nil {
			logger.Errorf("Failed to connect server %+v " + url.Location)
			return errors.New("Failed to connect server " + url.Location)
		}
	}
	//FIXME atomic operation
	cl.init = true
	return nil
}
