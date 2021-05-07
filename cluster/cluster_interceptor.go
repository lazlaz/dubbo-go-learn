package cluster

import (
	"context"
	"github.com/laz/dubbo-go/protocol"
)

// ClusterInterceptor
// Extension - ClusterInterceptor
type ClusterInterceptor interface {
	// Before DoInvoke method
	BeforeInvoker(ctx context.Context, invocation protocol.Invocation)

	// After DoInvoke method
	AfterInvoker(ctx context.Context, invocation protocol.Invocation)

	// Corresponding cluster invoke
	DoInvoke(ctx context.Context, invocation protocol.Invocation) protocol.Result
}
