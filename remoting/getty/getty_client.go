package getty

import (
	getty "github.com/apache/dubbo-getty"
	gxsync "github.com/dubbogo/gost/sync"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/config"
	"github.com/laz/dubbo-go/remoting"
	perrors "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"math/rand"
	"sync"
	"time"
)

var (
	errSessionNotExist   = perrors.New("session not exist")
	errClientClosed      = perrors.New("client closed")
	errClientReadTimeout = perrors.New("maybe the client read timeout or fail to decode tcp stream in Writer.Write")

	clientConf   *ClientConfig
	clientGrpool gxsync.GenericTaskPool
)

// Options : param config
type Options struct {
	// connect timeout
	// remove request timeout, it will be calculate for every request
	ConnectTimeout time.Duration
	// request timeout
	RequestTimeout time.Duration
}

// Client : some configuration for network communication.
type Client struct {
	addr           string
	opts           Options
	conf           ClientConfig
	mux            sync.RWMutex
	pool           *gettyRPCClientPool
	codec          remoting.Codec
	ExchangeClient *remoting.ExchangeClient
}

func (c *Client) SetExchangeClient(client *remoting.ExchangeClient) {
	c.ExchangeClient = client
}

// it is init client for single protocol.
func initClient(protocol string) {
	if protocol == "" {
		return
	}

	// load clientconfig from consumer_config
	// default use dubbo
	consumerConfig := config.GetConsumerConfig()
	if consumerConfig.ApplicationConfig == nil {
		return
	}
	protocolConf := config.GetConsumerConfig().ProtocolConf
	defaultClientConfig := GetDefaultClientConfig()
	if protocolConf == nil {
		logger.Info("protocol_conf default use dubbo config")
	} else {
		dubboConf := protocolConf.(map[interface{}]interface{})[protocol]
		if dubboConf == nil {
			logger.Warnf("dubboConf is nil")
			return
		}
		dubboConfByte, err := yaml.Marshal(dubboConf)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(dubboConfByte, &defaultClientConfig)
		if err != nil {
			panic(err)
		}
	}
	clientConf = &defaultClientConfig
	if err := clientConf.CheckValidity(); err != nil {
		logger.Warnf("[CheckValidity] error: %v", err)
		return
	}
	setClientGrpool()

	rand.Seed(time.Now().UnixNano())
}
func setClientGrpool() {
	clientGrpool = gxsync.NewTaskPoolSimple(clientConf.GrPoolSize)
}

func (c *Client) Connect(url *common.URL) error {
	initClient(url.Protocol)
	c.conf = *clientConf
	// new client
	c.pool = newGettyRPCClientConnPool(c, clientConf.PoolSize, time.Duration(int(time.Second)*clientConf.PoolTTL))
	c.pool.sslEnabled = url.GetParamBool(constant.SSL_ENABLED_KEY, false)

	// codec
	c.codec = remoting.GetCodec(url.Protocol)
	c.addr = url.Location
	_, _, err := c.selectSession(c.addr)
	if err != nil {
		logger.Errorf("try to connect server %v failed for : %v", url.Location, err)
	}
	return err
}

func (c *Client) selectSession(addr string) (*gettyRPCClient, getty.Session, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if c.pool == nil {
		return nil, nil, perrors.New("client pool have been closed")
	}
	rpcClient, err := c.pool.getGettyRpcClient(addr)
	if err != nil {
		return nil, nil, perrors.WithStack(err)
	}
	return rpcClient, rpcClient.selectSession(), nil
}
func (c *Client) Close() {
	c.mux.Lock()
	p := c.pool
	c.pool = nil
	c.mux.Unlock()
	if p != nil {
		p.close()
	}
}

func (c *Client) transfer(session getty.Session, request *remoting.Request, timeout time.Duration) (int, int, error) {
	totalLen, sendLen, err := session.WritePkg(request, timeout)
	return totalLen, sendLen, perrors.WithStack(err)
}

func (c *Client) Request(request *remoting.Request, timeout time.Duration, response *remoting.PendingResponse) error {
	_, session, err := c.selectSession(c.addr)
	if err != nil {
		return perrors.WithStack(err)
	}
	if session == nil {
		return errSessionNotExist
	}
	var (
		totalLen int
		sendLen  int
	)
	if totalLen, sendLen, err = c.transfer(session, request, timeout); err != nil {
		if sendLen != 0 && totalLen != sendLen {
			logger.Warnf("start to close the session at request because %d of %d bytes data is sent success. err:%+v", sendLen, totalLen, err)
			go c.Close()
		}
		return perrors.WithStack(err)
	}

	if !request.TwoWay || response.Callback != nil {
		return nil
	}

	select {
	case <-getty.GetTimeWheel().After(timeout):
		return perrors.WithStack(errClientReadTimeout)
	case <-response.Done:
		err = response.Err
	}

	return perrors.WithStack(err)
}

// isAvailable returns true if the connection is available, or it can be re-established.
func (c *Client) IsAvailable() bool {
	client, _, err := c.selectSession(c.addr)
	return err == nil &&
		// defensive check
		client != nil
}

// create client
func NewClient(opt Options) *Client {
	switch {
	case opt.ConnectTimeout == 0:
		opt.ConnectTimeout = 3 * time.Second
		fallthrough
	case opt.RequestTimeout == 0:
		opt.RequestTimeout = 3 * time.Second
	}

	c := &Client{
		opts: opt,
	}
	return c
}
