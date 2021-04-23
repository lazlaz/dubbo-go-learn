package getty

import (
	"crypto/tls"
	"fmt"
	getty "github.com/apache/dubbo-getty"
	_ "github.com/apache/dubbo-go-hessian2"
	gxsync "github.com/dubbogo/gost/sync"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/config"
	"github.com/laz/dubbo-go/protocol"
	"github.com/laz/dubbo-go/protocol/invocation"
	"github.com/laz/dubbo-go/remoting"
	perrors "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"net"
	"sync"
	"time"
)

// Server define getty server
type Server struct {
	conf           ServerConfig
	addr           string
	codec          remoting.Codec
	tcpServer      getty.Server
	rpcHandler     *RpcServerHandler
	requestHandler func(*invocation.RPCInvocation) protocol.RPCResult
}

// nolint
type RpcServerHandler struct {
	maxSessionNum  int
	sessionTimeout time.Duration
	sessionMap     map[getty.Session]*rpcSession
	rwlock         sync.RWMutex
	server         *Server
	timeoutTimes   int
}

// OnCron check the session health periodic. if the session's sessionTimeout has reached, just close the session
func (h *RpcServerHandler) OnCron(session getty.Session) {

}

// OnOpen call server session opened, add the session to getty server session list. also onOpen
// will check the max getty server session number
func (h *RpcServerHandler) OnOpen(session getty.Session) error {
	var err error
	h.rwlock.RLock()
	if h.maxSessionNum <= len(h.sessionMap) {
		err = errTooManySessions
	}
	h.rwlock.RUnlock()
	if err != nil {
		return perrors.WithStack(err)
	}

	logger.Infof("got session:%s", session.Stat())
	h.rwlock.Lock()
	h.sessionMap[session] = &rpcSession{session: session}
	h.rwlock.Unlock()
	return nil
}

// OnError the getty server session has errored, so remove the session from the getty server session list
func (h *RpcServerHandler) OnError(session getty.Session, err error) {
	logger.Infof("session{%s} got error{%v}, will be closed.", session.Stat(), err)
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

// OnClose close the session, remove it from the getty server list
func (h *RpcServerHandler) OnClose(session getty.Session) {
	logger.Infof("session{%s} is closing......", session.Stat())
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

// OnMessage get request from getty client, update the session reqNum and reply response to client
func (h *RpcServerHandler) OnMessage(session getty.Session, pkg interface{}) {

}

var (
	srvConf *ServerConfig
)

// NewServer create a new Server
func NewServer(url *common.URL, handlers func(*invocation.RPCInvocation) protocol.RPCResult) *Server {
	//init
	initServer(url.Protocol)

	srvConf.SSLEnabled = url.GetParamBool(constant.SSL_ENABLED_KEY, false)

	s := &Server{
		conf:           *srvConf,
		addr:           url.Location,
		codec:          remoting.GetCodec(url.Protocol),
		requestHandler: handlers,
	}

	s.rpcHandler = NewRpcServerHandler(s.conf.SessionNumber, s.conf.sessionTimeout, s)

	return s
}

// GetServerConfig get getty server config.
func GetServerConfig() ServerConfig {
	return *srvConf
}
func initServer(protocol string) {
	// load clientconfig from provider_config
	// default use dubbo
	providerConfig := config.GetProviderConfig()
	if providerConfig.ApplicationConfig == nil {
		return
	}
	protocolConf := providerConfig.ProtocolConf
	defaultServerConfig := GetDefaultServerConfig()
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
		err = yaml.Unmarshal(dubboConfByte, &defaultServerConfig)
		if err != nil {
			panic(err)
		}
	}
	srvConf = &defaultServerConfig
	if err := srvConf.CheckValidity(); err != nil {
		panic(err)
	}
}
func (s *Server) newSession(session getty.Session) error {
	var (
		ok      bool
		tcpConn *net.TCPConn
		err     error
	)
	conf := s.conf

	if conf.GettySessionParam.CompressEncoding {
		session.SetCompressType(getty.CompressZip)
	}
	if _, ok = session.Conn().(*tls.Conn); ok {
		session.SetName(conf.GettySessionParam.SessionName)
		session.SetMaxMsgLen(conf.GettySessionParam.MaxMsgLen)
		session.SetPkgHandler(NewRpcServerPackageHandler(s))
		session.SetEventListener(s.rpcHandler)
		session.SetReadTimeout(conf.GettySessionParam.tcpReadTimeout)
		session.SetWriteTimeout(conf.GettySessionParam.tcpWriteTimeout)
		session.SetCronPeriod((int)(conf.heartbeatPeriod.Nanoseconds() / 1e6))
		session.SetWaitTime(conf.GettySessionParam.waitTimeout)
		logger.Debugf("server accepts new session:%s\n", session.Stat())
		return nil
	}
	if _, ok = session.Conn().(*net.TCPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
	}

	if _, ok = session.Conn().(*tls.Conn); !ok {
		if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
			return perrors.New(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection", session.Stat(), session.Conn()))
		}

		if err = tcpConn.SetNoDelay(conf.GettySessionParam.TcpNoDelay); err != nil {
			return err
		}
		if err = tcpConn.SetKeepAlive(conf.GettySessionParam.TcpKeepAlive); err != nil {
			return err
		}
		if conf.GettySessionParam.TcpKeepAlive {
			if err = tcpConn.SetKeepAlivePeriod(conf.GettySessionParam.keepAlivePeriod); err != nil {
				return err
			}
		}
		if err = tcpConn.SetReadBuffer(conf.GettySessionParam.TcpRBufSize); err != nil {
			return err
		}
		if err = tcpConn.SetWriteBuffer(conf.GettySessionParam.TcpWBufSize); err != nil {
			return err
		}
	}

	session.SetName(conf.GettySessionParam.SessionName)
	session.SetMaxMsgLen(conf.GettySessionParam.MaxMsgLen)
	session.SetPkgHandler(NewRpcServerPackageHandler(s))
	session.SetEventListener(s.rpcHandler)
	session.SetReadTimeout(conf.GettySessionParam.tcpReadTimeout)
	session.SetWriteTimeout(conf.GettySessionParam.tcpWriteTimeout)
	session.SetCronPeriod((int)(conf.heartbeatPeriod.Nanoseconds() / 1e6))
	session.SetWaitTime(conf.GettySessionParam.waitTimeout)
	logger.Debugf("server accepts new session: %s", session.Stat())
	return nil
}

// Start dubbo server.
func (s *Server) Start() {
	var (
		addr      string
		tcpServer getty.Server
	)

	addr = s.addr
	serverOpts := []getty.ServerOption{getty.WithLocalAddress(addr)}
	if s.conf.SSLEnabled {
		serverOpts = append(serverOpts, getty.WithServerSslEnabled(s.conf.SSLEnabled),
			getty.WithServerTlsConfigBuilder(config.GetServerTlsConfigBuilder()))
	}

	serverOpts = append(serverOpts, getty.WithServerTaskPool(gxsync.NewTaskPoolSimple(s.conf.GrPoolSize)))

	tcpServer = getty.NewTCPServer(serverOpts...)
	tcpServer.RunEventLoop(s.newSession)
	logger.Debugf("s bind addr{%s} ok!", s.addr)
	s.tcpServer = tcpServer
}

// Stop dubbo server
func (s *Server) Stop() {
	s.tcpServer.Close()
}
