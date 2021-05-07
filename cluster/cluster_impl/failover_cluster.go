package cluster_impl

import (
	"github.com/laz/dubbo-go/cluster"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/extension"
	"github.com/laz/dubbo-go/protocol"
)

type failoverCluster struct{}

func init() {
	extension.SetCluster(constant.FAILOVER_CLUSTER_NAME, NewFailoverCluster)
}
func NewFailoverCluster() cluster.Cluster {
	return &failoverCluster{}
}

// Join returns a baseClusterInvoker instance
func (cluster *failoverCluster) Join(directory cluster.Directory) protocol.Invoker {
	return newFailoverClusterInvoker(directory)
}
