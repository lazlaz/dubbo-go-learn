package getty

import (
	"crypto/tls"
	"fmt"
	"github.com/laz/dubbo-go/common/logger"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)
import (
	"github.com/apache/dubbo-getty"
	perrors "github.com/pkg/errors"
)

type gettyRPCClient struct {
	once sync.Once
	//protocol string
	addr   string
	active int64 // zero, not create or be destroyed

	pool *gettyRPCClientPool

	lock        sync.RWMutex
	gettyClient getty.Client
	sessions    []*rpcSession
}
type gettyRPCClientPool struct {
	rpcClient  *Client
	size       int   // size of []*gettyRPCClient
	ttl        int64 // ttl of every gettyRPCClient, it is checked when getConn
	sslEnabled bool

	sync.Mutex
	conns []*gettyRPCClient
}

func newGettyRPCClientConnPool(rpcClient *Client, size int, ttl time.Duration) *gettyRPCClientPool {
	return &gettyRPCClientPool{
		rpcClient: rpcClient,
		size:      size,
		ttl:       int64(ttl.Seconds()),
		// init capacity : 2
		conns: make([]*gettyRPCClient, 0, 2),
	}
}

func (p *gettyRPCClientPool) getGettyRpcClient(addr string) (*gettyRPCClient, error) {
	conn, connErr := p.get()
	if connErr == nil && conn == nil {
		// create new conn
		rpcClientConn, rpcErr := newGettyRPCClientConn(p, addr)
		if rpcErr == nil {
			p.put(rpcClientConn)
		}
		return rpcClientConn, perrors.WithStack(rpcErr)
	}
	return conn, perrors.WithStack(connErr)
}
func (c *gettyRPCClient) getActive() int64 {
	return atomic.LoadInt64(&c.active)
}
func (p *gettyRPCClientPool) put(conn *gettyRPCClient) {
	if conn == nil || conn.getActive() == 0 {
		return
	}
	p.Lock()
	defer p.Unlock()
	if p.conns == nil {
		return
	}
	// check whether @conn has existed in p.conns or not.
	for i := range p.conns {
		if p.conns[i] == conn {
			return
		}
	}
	if len(p.conns) >= p.size {
		// delete @conn from client pool
		// p.remove(conn)
		conn.close()
		return
	}
	p.conns = append(p.conns, conn)
}
func (p *gettyRPCClientPool) close() {
	p.Lock()
	conns := p.conns
	p.conns = nil
	p.Unlock()
	for _, conn := range conns {
		conn.close()
	}
}

func (c *gettyRPCClient) close() error {
	closeErr := perrors.Errorf("close gettyRPCClient{%#v} again", c)
	c.once.Do(func() {
		var (
			gettyClient getty.Client
			sessions    []*rpcSession
		)
		func() {
			c.lock.Lock()
			defer c.lock.Unlock()

			gettyClient = c.gettyClient
			c.gettyClient = nil

			sessions = make([]*rpcSession, 0, len(c.sessions))
			sessions = append(sessions, c.sessions...)
			c.sessions = c.sessions[:0]
		}()

		c.updateActive(0)

		go func() {
			if gettyClient != nil {
				gettyClient.Close()
			}
			for _, s := range sessions {
				logger.Infof("close client session{%s, last active:%s, request number:%d}",
					s.session.Stat(), s.session.GetActive().String(), s.GetReqNum())
				s.session.Close()
			}
		}()

		closeErr = nil
	})
	return closeErr
}
func (c *gettyRPCClient) updateActive(active int64) {
	atomic.StoreInt64(&c.active, active)
}
func (p *gettyRPCClientPool) remove(conn *gettyRPCClient) {
	if conn == nil || conn.getActive() == 0 {
		return
	}

	if p.conns == nil {
		return
	}

	if len(p.conns) > 0 {
		for idx, c := range p.conns {
			if conn == c {
				p.conns = append(p.conns[:idx], p.conns[idx+1:]...)
				break
			}
		}
	}
}
func newGettyRPCClientConn(pool *gettyRPCClientPool, addr string) (*gettyRPCClient, error) {
	var (
		gettyClient getty.Client
	)

	clientOpts := []getty.ClientOption{
		getty.WithServerAddress(addr),
		getty.WithConnectionNumber((int)(pool.rpcClient.conf.ConnectionNum)),
		getty.WithReconnectInterval(pool.rpcClient.conf.ReconnectInterval),
	}

	if clientGrpool != nil {
		clientOpts = append(clientOpts, getty.WithClientTaskPool(clientGrpool))
	}

	gettyClient = getty.NewTCPClient(clientOpts...)
	c := &gettyRPCClient{
		addr:        addr,
		pool:        pool,
		gettyClient: gettyClient,
	}
	go c.gettyClient.RunEventLoop(c.newSession)

	idx := 1
	start := time.Now()
	connectTimeout := pool.rpcClient.opts.ConnectTimeout
	for {
		idx++
		if c.isAvailable() {
			break
		}

		if time.Since(start) > connectTimeout {
			c.gettyClient.Close()
			return nil, perrors.New(fmt.Sprintf("failed to create client connection to %s in %s", addr, connectTimeout))
		}

		interval := time.Millisecond * time.Duration(idx)
		if interval > time.Duration(100e6) {
			interval = 100e6 // 100 ms
		}
		time.Sleep(interval)
	}
	logger.Debug("client init ok")
	c.updateActive(time.Now().Unix())

	return c, nil
}
func (c *gettyRPCClient) isAvailable() bool {
	return c.selectSession() != nil
}
func (c *gettyRPCClient) newSession(session getty.Session) error {
	var (
		ok         bool
		tcpConn    *net.TCPConn
		conf       ClientConfig
		sslEnabled bool
	)
	conf = c.pool.rpcClient.conf
	sslEnabled = c.pool.sslEnabled
	if conf.GettySessionParam.CompressEncoding {
		session.SetCompressType(getty.CompressZip)
	}
	if sslEnabled {
		if _, ok = session.Conn().(*tls.Conn); !ok {
			panic(fmt.Sprintf("%s, session.conn{%#v} is not tls connection\n", session.Stat(), session.Conn()))
		}
		session.SetName(conf.GettySessionParam.SessionName)
		session.SetMaxMsgLen(conf.GettySessionParam.MaxMsgLen)
		session.SetPkgHandler(NewRpcClientPackageHandler(c.pool.rpcClient))
		session.SetEventListener(NewRpcClientHandler(c))
		session.SetReadTimeout(conf.GettySessionParam.tcpReadTimeout)
		session.SetWriteTimeout(conf.GettySessionParam.tcpWriteTimeout)
		session.SetCronPeriod((int)(conf.heartbeatPeriod.Nanoseconds() / 1e6))
		session.SetWaitTime(conf.GettySessionParam.waitTimeout)
		logger.Debugf("client new session:%s\n", session.Stat())
		return nil
	}
	if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
	}

	if err := tcpConn.SetNoDelay(conf.GettySessionParam.TcpNoDelay); err != nil {
		logger.Error("tcpConn.SetNoDelay() = error:%v", err)
	}
	if err := tcpConn.SetKeepAlive(conf.GettySessionParam.TcpKeepAlive); err != nil {
		logger.Error("tcpConn.SetKeepAlive() = error:%v", err)
	}
	if conf.GettySessionParam.TcpKeepAlive {
		if err := tcpConn.SetKeepAlivePeriod(conf.GettySessionParam.keepAlivePeriod); err != nil {
			logger.Error("tcpConn.SetKeepAlivePeriod() = error:%v", err)
		}
	}
	if err := tcpConn.SetReadBuffer(conf.GettySessionParam.TcpRBufSize); err != nil {
		logger.Error("tcpConn.SetReadBuffer() = error:%v", err)
	}
	if err := tcpConn.SetWriteBuffer(conf.GettySessionParam.TcpWBufSize); err != nil {
		logger.Error("tcpConn.SetWriteBuffer() = error:%v", err)
	}

	session.SetName(conf.GettySessionParam.SessionName)
	session.SetMaxMsgLen(conf.GettySessionParam.MaxMsgLen)
	session.SetPkgHandler(NewRpcClientPackageHandler(c.pool.rpcClient))
	session.SetEventListener(NewRpcClientHandler(c))
	session.SetReadTimeout(conf.GettySessionParam.tcpReadTimeout)
	session.SetWriteTimeout(conf.GettySessionParam.tcpWriteTimeout)
	session.SetCronPeriod((int)(conf.heartbeatPeriod.Nanoseconds() / 1e6))
	session.SetWaitTime(conf.GettySessionParam.waitTimeout)
	logger.Debugf("client new session:%s\n", session.Stat())
	return nil
}

func (c *gettyRPCClient) addSession(session getty.Session) {
	logger.Debugf("add session{%s}", session.Stat())
	if session == nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	if c.sessions == nil {
		c.sessions = make([]*rpcSession, 0, 16)
	}
	c.sessions = append(c.sessions, &rpcSession{session: session})
}
func (c *gettyRPCClient) updateSession(session getty.Session) {
	if session == nil {
		return
	}

	var rs *rpcSession
	func() {
		c.lock.RLock()
		defer c.lock.RUnlock()
		if c.sessions == nil {
			return
		}

		for i, s := range c.sessions {
			if s.session == session {
				rs = c.sessions[i]
				break
			}
		}
	}()
	if rs != nil {
		rs.AddReqNum(1)
	}
}

func (c *gettyRPCClient) selectSession() getty.Session {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.sessions == nil {
		return nil
	}

	count := len(c.sessions)
	if count == 0 {
		return nil
	}
	return c.sessions[rand.Int31n(int32(count))].session
}

var (
	errClientPoolClosed = perrors.New("client pool closed")
)

func (p *gettyRPCClientPool) get() (*gettyRPCClient, error) {
	now := time.Now().Unix()

	p.Lock()
	defer p.Unlock()
	if p.conns == nil {
		return nil, errClientPoolClosed
	}
	for num := len(p.conns); num > 0; {
		var conn *gettyRPCClient
		if num != 1 {
			conn = p.conns[rand.Int31n(int32(num))]
		} else {
			conn = p.conns[0]
		}
		// This will recreate gettyRpcClient for remove last position
		//p.conns = p.conns[:len(p.conns)-1]

		if d := now - conn.getActive(); d > p.ttl {
			p.remove(conn)
			go conn.close()
			num = len(p.conns)
			continue
		}
		conn.updateActive(now) //update active time
		return conn, nil
	}
	return nil, nil
}
