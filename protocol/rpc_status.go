package protocol

import (
	"github.com/laz/dubbo-go/common/logger"
	"sync"
)

import (
	uberAtomic "go.uber.org/atomic"
)

var (
	methodStatistics    sync.Map        // url -> { methodName : RPCStatus}
	serviceStatistic    sync.Map        // url -> RPCStatus
	invokerBlackList    sync.Map        // store unhealthy url blackList
	blackListCacheDirty uberAtomic.Bool // store if the cache in chain is not refreshed by blacklist
	blackListRefreshing int32           // store if the refresing method is processing
)

// RemoveUrlKeyUnhealthyStatus called when event of provider unregister, delete from black list
func RemoveUrlKeyUnhealthyStatus(key string) {
	invokerBlackList.Delete(key)
	logger.Info("Remove invoker key = ", key, " from black list")
	blackListCacheDirty.Store(true)
}
