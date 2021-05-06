package extension

import "github.com/laz/dubbo-go/cluster"

var (
	clusters = make(map[string]func() cluster.Cluster)
)

// SetCluster sets the cluster fault-tolerant mode with @name
// For example: available/failfast/broadcast/failfast/failsafe/...
func SetCluster(name string, fcn func() cluster.Cluster) {
	clusters[name] = fcn
}

// GetCluster finds the cluster fault-tolerant mode with @name
func GetCluster(name string) cluster.Cluster {
	if clusters[name] == nil {
		panic("cluster for " + name + " is not existing, make sure you have import the package.")
	}
	return clusters[name]()
}
