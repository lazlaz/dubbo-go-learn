package cluster_impl

import (
	"context"
	"github.com/laz/dubbo-go/cluster"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/protocol"
	"go.uber.org/atomic"
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
