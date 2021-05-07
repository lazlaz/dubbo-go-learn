package getty

import (
	getty "github.com/apache/dubbo-getty"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/remoting"
	perrors "github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type rpcSession struct {
	session getty.Session
	reqNum  int32
}

// todo: WritePkg_Timeout will entry *.yml
const (
	// WritePkg_Timeout the timeout of write pkg
	WritePkg_Timeout = 5 * time.Second
)

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

// nolint
type RpcClientHandler struct {
	conn         *gettyRPCClient
	timeoutTimes int
}

// nolint
func NewRpcClientHandler(client *gettyRPCClient) *RpcClientHandler {
	return &RpcClientHandler{conn: client}
}

// OnOpen call the getty client session opened, add the session to getty client session list
func (h *RpcClientHandler) OnOpen(session getty.Session) error {
	h.conn.addSession(session)
	return nil
}

// OnError the getty client session has errored, so remove the session from the getty client session list
func (h *RpcClientHandler) OnError(session getty.Session, err error) {
	logger.Infof("session{%s} got error{%v}, will be closed.", session.Stat(), err)

}

// OnClose close the session, remove it from the getty session list
func (h *RpcClientHandler) OnClose(session getty.Session) {
	logger.Infof("session{%s} is closing......", session.Stat())

}

// OnMessage get response from getty server, and update the session to the getty client session list
func (h *RpcClientHandler) OnMessage(session getty.Session, pkg interface{}) {
	result, ok := pkg.(remoting.DecodeResult)
	if !ok {
		logger.Errorf("illegal package")
		return
	}
	// get heartbeat request from server
	if result.IsRequest {
		req := result.Result.(*remoting.Request)
		if req.Event {
			logger.Debugf("get rpc heartbeat request{%#v}", req)
			resp := remoting.NewResponse(req.ID, req.Version)
			resp.Status = hessian.Response_OK
			resp.Event = req.Event
			resp.SerialID = req.SerialID
			resp.Version = "2.0.2"
			reply(session, resp)
			return
		}
		logger.Errorf("illegal request but not heartbeat. {%#v}", req)
		return
	}
	h.timeoutTimes = 0
	p := result.Result.(*remoting.Response)
	// get heartbeat
	if p.Event {
		logger.Debugf("get rpc heartbeat response{%#v}", p)
		if p.Error != nil {
			logger.Errorf("rpc heartbeat response{error: %#v}", p.Error)
		}
		p.Handle()
		return
	}

	logger.Debugf("get rpc response{%#v}", p)

	h.conn.updateSession(session)

	p.Handle()
}

// OnCron check the session health periodic. if the session's sessionTimeout has reached, just close the session
func (h *RpcClientHandler) OnCron(session getty.Session) {

}
func reply(session getty.Session, resp *remoting.Response) {
	if totalLen, sendLen, err := session.WritePkg(resp, WritePkg_Timeout); err != nil {
		if sendLen != 0 && totalLen != sendLen {
			logger.Warnf("start to close the session at replying because %d of %d bytes data is sent success. err:%+v", sendLen, totalLen, err)
			go session.Close()
		}
		logger.Errorf("WritePkg error: %#v, %#v", perrors.WithStack(err), resp)
	}
}
