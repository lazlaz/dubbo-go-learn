package constant

const (
	DEFAULT_KEY               = "default"
	PREFIX_DEFAULT_KEY        = "default."
	DEFAULT_SERVICE_FILTERS   = ""
	DEFAULT_REFERENCE_FILTERS = "cshutdown"
	GENERIC_REFERENCE_FILTERS = "generic"
	GENERIC                   = "$invoke"
	ECHO                      = "$echo"
)
const (
	DEFAULT_WEIGHT = 100     //
	DEFAULT_WARMUP = 10 * 60 // in java here is 10*60*1000 because of System.currentTimeMillis() is measured in milliseconds & in go time.Unix() is second
)
const (
	CONFIGURATORS_CATEGORY             = "configurators"
	ROUTER_CATEGORY                    = "category"
	DEFAULT_CATEGORY                   = PROVIDER_CATEGORY
	DYNAMIC_CONFIGURATORS_CATEGORY     = "dynamicconfigurators"
	APP_DYNAMIC_CONFIGURATORS_CATEGORY = "appdynamicconfigurators"
	PROVIDER_CATEGORY                  = "providers"
	CONSUMER_CATEGORY                  = "consumers"
)
const (
	DEFAULT_LOADBALANCE        = "random"
	DEFAULT_RETRIES            = "2"
	DEFAULT_RETRIES_INT        = 2
	DEFAULT_PROTOCOL           = "dubbo"
	DEFAULT_REG_TIMEOUT        = "10s"
	DEFAULT_REG_TTL            = "15m"
	DEFAULT_CLUSTER            = "failover"
	DEFAULT_FAILBACK_TIMES     = "3"
	DEFAULT_FAILBACK_TIMES_INT = 3
	DEFAULT_FAILBACK_TASKS     = 100
	DEFAULT_REST_CLIENT        = "resty"
	DEFAULT_REST_SERVER        = "go-restful"
	DEFAULT_PORT               = 20000
	DEFAULT_SERIALIZATION      = HESSIAN2_SERIALIZATION
)

const (
	DUBBO             = "dubbo"
	PROVIDER_PROTOCOL = "provider"
	//compatible with 2.6.x
	OVERRIDE_PROTOCOL = "override"
	EMPTY_PROTOCOL    = "empty"
	ROUTER_PROTOCOL   = "router"
)
