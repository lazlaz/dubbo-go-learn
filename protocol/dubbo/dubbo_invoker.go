package dubbo

import (
	"context"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/config"
	"github.com/laz/dubbo-go/protocol"
	invocation_impl "github.com/laz/dubbo-go/protocol/invocation"
	"github.com/laz/dubbo-go/remoting"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	attachmentKey = []string{constant.INTERFACE_KEY, constant.GROUP_KEY, constant.TOKEN_KEY, constant.TIMEOUT_KEY,
		constant.VERSION_KEY}
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

// get timeout including methodConfig
func (di *DubboInvoker) getTimeout(invocation *invocation_impl.RPCInvocation) time.Duration {
	var timeout = di.GetUrl().GetParam(strings.Join([]string{constant.METHOD_KEYS, invocation.MethodName(), constant.TIMEOUT_KEY}, "."), "")
	if len(timeout) != 0 {
		if t, err := time.ParseDuration(timeout); err == nil {
			// config timeout into attachment
			invocation.SetAttachments(constant.TIMEOUT_KEY, strconv.Itoa(int(t.Milliseconds())))
			return t
		}
	}
	// set timeout into invocation at method level
	invocation.SetAttachments(constant.TIMEOUT_KEY, strconv.Itoa(int(di.timeout.Milliseconds())))
	return di.timeout
}

// Invoke call remoting.
func (di *DubboInvoker) Invoke(ctx context.Context, invocation protocol.Invocation) protocol.Result {
	var (
		err    error
		result protocol.RPCResult
	)
	if !di.BaseInvoker.IsAvailable() {
		// Generally, the case will not happen, because the invoker has been removed
		// from the invoker list before destroy,so no new request will enter the destroyed invoker
		logger.Warnf("this dubboInvoker is destroyed")
		result.Err = protocol.ErrDestroyedInvoker
		return &result
	}

	di.clientGuard.RLock()
	defer di.clientGuard.RUnlock()

	if di.client == nil {
		result.Err = protocol.ErrClientClosed
		logger.Debugf("result.Err: %v", result.Err)
		return &result
	}

	if !di.BaseInvoker.IsAvailable() {
		// Generally, the case will not happen, because the invoker has been removed
		// from the invoker list before destroy,so no new request will enter the destroyed invoker
		logger.Warnf("this dubboInvoker is destroying")
		result.Err = protocol.ErrDestroyedInvoker
		return &result
	}

	inv := invocation.(*invocation_impl.RPCInvocation)
	// init param
	inv.SetAttachments(constant.PATH_KEY, di.GetUrl().GetParam(constant.INTERFACE_KEY, ""))
	for _, k := range attachmentKey {
		if v := di.GetUrl().GetParam(k, ""); len(v) > 0 {
			inv.SetAttachments(k, v)
		}
	}

	// put the ctx into attachment
	di.appendCtx(ctx, inv)

	url := di.GetUrl()
	// default hessian2 serialization, compatible
	if url.GetParam(constant.SERIALIZATION_KEY, "") == "" {
		url.SetParam(constant.SERIALIZATION_KEY, constant.HESSIAN2_SERIALIZATION)
	}
	// async
	async, err := strconv.ParseBool(inv.AttachmentsByKey(constant.ASYNC_KEY, "false"))
	if err != nil {
		logger.Errorf("ParseBool - error: %v", err)
		async = false
	}
	//response := NewResponse(inv.Reply(), nil)
	rest := &protocol.RPCResult{}
	timeout := di.getTimeout(inv)
	if async {
		if callBack, ok := inv.CallBack().(func(response common.CallbackResponse)); ok {
			result.Err = di.client.AsyncRequest(&invocation, url, timeout, callBack, rest)
		} else {
			result.Err = di.client.Send(&invocation, url, timeout)
		}
	} else {
		if inv.Reply() == nil {
			result.Err = protocol.ErrNoReply
		} else {
			result.Err = di.client.Request(&invocation, url, timeout, rest)
		}
	}
	if result.Err == nil {
		result.Rest = inv.Reply()
		result.Attrs = rest.Attrs
	}
	logger.Debugf("result.Err: %v, result.Rest: %v", result.Err, result.Rest)

	return &result
}

// Finally, I made the decision that I don't provide a general way to transfer the whole context
// because it could be misused. If the context contains to many key-value pairs, the performance will be much lower.
func (di *DubboInvoker) appendCtx(ctx context.Context, inv *invocation_impl.RPCInvocation) {
	// inject opentracing ctx
	/*	currentSpan := opentracing.SpanFromContext(ctx)
		if currentSpan != nil {
			err := injectTraceCtx(currentSpan, inv)
			if err != nil {
				logger.Errorf("Could not inject the span context into attachments: %v", err)
			}
		}*/
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
