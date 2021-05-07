package extension

import "github.com/laz/dubbo-go/cluster"

var (
	loadbalances = make(map[string]func() cluster.LoadBalance)
)

// SetLoadbalance sets the loadbalance extension with @name
// For example: random/round_robin/consistent_hash/least_active/...
func SetLoadbalance(name string, fcn func() cluster.LoadBalance) {
	loadbalances[name] = fcn
}

// GetLoadbalance finds the loadbalance extension with @name
func GetLoadbalance(name string) cluster.LoadBalance {
	if loadbalances[name] == nil {
		panic("loadbalance for " + name + " is not existing, make sure you have import the package.")
	}

	return loadbalances[name]()
}
