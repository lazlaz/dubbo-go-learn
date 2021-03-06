package cluster_impl

import (
	"context"
	"github.com/laz/dubbo-go/cluster"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/protocol"
	"go.uber.org/atomic"
)
import (
	perrors "github.com/pkg/errors"
)

type baseClusterInvoker struct {
	directory      cluster.Directory
	availablecheck bool
	destroyed      *atomic.Bool
	stickyInvoker  protocol.Invoker
	interceptor    cluster.ClusterInterceptor
}

func newBaseClusterInvoker(directory cluster.Directory) baseClusterInvoker {
	return baseClusterInvoker{
		directory:      directory,
		availablecheck: true,
		destroyed:      atomic.NewBool(false),
	}
}

func (invoker *baseClusterInvoker) GetUrl() *common.URL {
	return invoker.directory.GetUrl()
}

func (invoker *baseClusterInvoker) Destroy() {
	//this is must atom operation
	if invoker.destroyed.CAS(false, true) {
		invoker.directory.Destroy()
	}
}
func (invoker *baseClusterInvoker) IsAvailable() bool {
	if invoker.stickyInvoker != nil {
		return invoker.stickyInvoker.IsAvailable()
	}
	return invoker.directory.IsAvailable()
}
func (invoker *baseClusterInvoker) Invoke(ctx context.Context, invocation protocol.Invocation) protocol.Result {
	if invoker.interceptor != nil {
		invoker.interceptor.BeforeInvoker(ctx, invocation)

		result := invoker.interceptor.DoInvoke(ctx, invocation)

		invoker.interceptor.AfterInvoker(ctx, invocation)

		return result
	}

	return nil
}

//check invokers availables
func (invoker *baseClusterInvoker) checkInvokers(invokers []protocol.Invoker, invocation protocol.Invocation) error {
	if len(invokers) == 0 {
		ip := common.GetLocalIp()
		return perrors.Errorf("Failed to invoke the method %v. No provider available for the service %v from "+
			"registry %v on the consumer %v using the dubbo version %v .Please check if the providers have been started and registered.",
			invocation.MethodName(), invoker.directory.GetUrl().SubURL.Key(), invoker.directory.GetUrl().String(), ip, constant.Version)
	}
	return nil

}
func getLoadBalance(invoker protocol.Invoker, invocation protocol.Invocation) cluster.LoadBalance {
	url := invoker.GetUrl()

	methodName := invocation.MethodName()
	//Get the service loadbalance config
	lb := url.GetParam(constant.LOADBALANCE_KEY, constant.DEFAULT_LOADBALANCE)

	//Get the service method loadbalance config if have
	if v := url.GetMethodParam(methodName, constant.LOADBALANCE_KEY, ""); len(v) > 0 {
		lb = v
	}
	return extension.GetLoadbalance(lb)
}

//check cluster invoker is destroyed or not
func (invoker *baseClusterInvoker) checkWhetherDestroyed() error {
	if invoker.destroyed.Load() {
		ip := common.GetLocalIp()
		return perrors.Errorf("Rpc cluster invoker for %v on consumer %v use dubbo version %v is now destroyed! can not invoke any more. ",
			invoker.directory.GetUrl().Service(), ip, constant.Version)
	}
	return nil
}

func (invoker *baseClusterInvoker) doSelect(lb cluster.LoadBalance, invocation protocol.Invocation, invokers []protocol.Invoker, invoked []protocol.Invoker) protocol.Invoker {
	var selectedInvoker protocol.Invoker
	if len(invokers) <= 0 {
		return selectedInvoker
	}

	url := invokers[0].GetUrl()
	sticky := url.GetParamBool(constant.STICKY_KEY, false)
	//Get the service method sticky config if have
	sticky = url.GetMethodParamBool(invocation.MethodName(), constant.STICKY_KEY, sticky)

	if invoker.stickyInvoker != nil && !isInvoked(invoker.stickyInvoker, invokers) {
		invoker.stickyInvoker = nil
	}

	if sticky && invoker.availablecheck &&
		invoker.stickyInvoker != nil && invoker.stickyInvoker.IsAvailable() &&
		(invoked == nil || !isInvoked(invoker.stickyInvoker, invoked)) {
		return invoker.stickyInvoker
	}

	selectedInvoker = invoker.doSelectInvoker(lb, invocation, invokers, invoked)
	if sticky {
		invoker.stickyInvoker = selectedInvoker
	}
	return selectedInvoker
}
func isInvoked(selectedInvoker protocol.Invoker, invoked []protocol.Invoker) bool {
	for _, i := range invoked {
		if i == selectedInvoker {
			return true
		}
	}
	return false
}

func (invoker *baseClusterInvoker) doSelectInvoker(lb cluster.LoadBalance, invocation protocol.Invocation, invokers []protocol.Invoker, invoked []protocol.Invoker) protocol.Invoker {
	if len(invokers) == 0 {
		return nil
	}

	selectedInvoker := lb.Select(invokers, invocation)
	return selectedInvoker

}
