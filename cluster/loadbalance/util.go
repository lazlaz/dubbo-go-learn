package loadbalance

import (
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/protocol"
	"time"
)

// GetWeight gets weight for load balance strategy
func GetWeight(invoker protocol.Invoker, invocation protocol.Invocation) int64 {
	var weight int64
	url := invoker.GetUrl()
	// Multiple registry scenario, load balance among multiple registries.
	isRegIvk := url.GetParamBool(constant.REGISTRY_KEY+"."+constant.REGISTRY_LABEL_KEY, false)
	if isRegIvk {
		weight = url.GetParamInt(constant.REGISTRY_KEY+"."+constant.WEIGHT_KEY, constant.DEFAULT_WEIGHT)
	} else {
		weight = url.GetMethodParamInt64(invocation.MethodName(), constant.WEIGHT_KEY, constant.DEFAULT_WEIGHT)

		if weight > 0 {
			//get service register time an do warm up time
			now := time.Now().Unix()
			timestamp := url.GetParamInt(constant.REMOTE_TIMESTAMP_KEY, now)
			if uptime := now - timestamp; uptime > 0 {
				warmup := url.GetParamInt(constant.WARMUP_KEY, constant.DEFAULT_WARMUP)
				if uptime < warmup {
					if ww := float64(uptime) / float64(warmup) / float64(weight); ww < 1 {
						weight = 1
					} else if int64(ww) <= weight {
						weight = int64(ww)
					}
				}
			}
		}
	}

	if weight < 0 {
		weight = 0
	}

	return weight
}
