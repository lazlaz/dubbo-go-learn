package cluster_impl

import (
	"github.com/laz/dubbo-go/cluster"
	"github.com/laz/dubbo-go/protocol"
)

type failoverClusterInvoker struct {
	baseClusterInvoker
}

func newFailoverClusterInvoker(directory cluster.Directory) protocol.Invoker {
	return &failoverClusterInvoker{
		baseClusterInvoker: newBaseClusterInvoker(directory),
	}
}
