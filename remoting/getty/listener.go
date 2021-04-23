package getty

import (
	getty "github.com/apache/dubbo-getty"
	perrors "github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type rpcSession struct {
	session getty.Session
	reqNum  int32
}

var (
	errTooManySessions      = perrors.New("too many sessions")
	errHeartbeatReadTimeout = perrors.New("heartbeat read timeout")
)

func (s *rpcSession) AddReqNum(num int32) {
	atomic.AddInt32(&s.reqNum, num)
}

func (s *rpcSession) GetReqNum() int32 {
	return atomic.LoadInt32(&s.reqNum)
}

// nolint
func NewRpcServerHandler(maxSessionNum int, sessionTimeout time.Duration, serverP *Server) *RpcServerHandler {
	return &RpcServerHandler{
		maxSessionNum:  maxSessionNum,
		sessionTimeout: sessionTimeout,
		sessionMap:     make(map[getty.Session]*rpcSession),
		server:         serverP,
	}
}
