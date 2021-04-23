package constant

type DubboCtxKey string

const (
	GROUP_KEY                = "group"
	VERSION_KEY              = "version"
	INTERFACE_KEY            = "interface"
	MESSAGE_SIZE_KEY         = "message_size"
	PATH_KEY                 = "path"
	SERVICE_KEY              = "service"
	METHODS_KEY              = "methods"
	TIMEOUT_KEY              = "timeout"
	CATEGORY_KEY             = "category"
	CHECK_KEY                = "check"
	ENABLED_KEY              = "enabled"
	SIDE_KEY                 = "side"
	OVERRIDE_PROVIDERS_KEY   = "providerAddresses"
	BEAN_NAME_KEY            = "bean.name"
	GENERIC_KEY              = "generic"
	CLASSIFIER_KEY           = "classifier"
	TOKEN_KEY                = "token"
	LOCAL_ADDR               = "local-addr"
	REMOTE_ADDR              = "remote-addr"
	DEFAULT_REMOTING_TIMEOUT = 3000
	RELEASE_KEY              = "release"
	ANYHOST_KEY              = "anyhost"
	PORT_KEY                 = "port"
	PROTOCOL_KEY             = "protocol"
	PATH_SEPARATOR           = "/"
	//DUBBO_KEY                = "dubbo"
	SSL_ENABLED_KEY = "ssl-enabled"
)
const (
	REGISTRY_KEY         = "registry"
	REGISTRY_PROTOCOL    = "registry"
	ROLE_KEY             = "registry.role"
	REGISTRY_DEFAULT_KEY = "registry.default"
	REGISTRY_TIMEOUT_KEY = "registry.timeout"
	REGISTRY_LABEL_KEY   = "label"
	PREFERRED_KEY        = "preferred"
	ZONE_KEY             = "zone"
	ZONE_FORCE_KEY       = "zone.force"
	REGISTRY_TTL_KEY     = "registry.ttl"
)
const (
	TIMESTAMP_KEY                          = "timestamp"
	REMOTE_TIMESTAMP_KEY                   = "remote.timestamp"
	CLUSTER_KEY                            = "cluster"
	LOADBALANCE_KEY                        = "loadbalance"
	WEIGHT_KEY                             = "weight"
	WARMUP_KEY                             = "warmup"
	RETRIES_KEY                            = "retries"
	STICKY_KEY                             = "sticky"
	BEAN_NAME                              = "bean.name"
	FAIL_BACK_TASKS_KEY                    = "failbacktasks"
	FORKS_KEY                              = "forks"
	DEFAULT_FORKS                          = 2
	DEFAULT_TIMEOUT                        = 1000
	ACCESS_LOG_KEY                         = "accesslog"
	TPS_LIMITER_KEY                        = "tps.limiter"
	TPS_REJECTED_EXECUTION_HANDLER_KEY     = "tps.limit.rejected.handler"
	TPS_LIMIT_RATE_KEY                     = "tps.limit.rate"
	DEFAULT_TPS_LIMIT_RATE                 = "-1"
	TPS_LIMIT_INTERVAL_KEY                 = "tps.limit.interval"
	DEFAULT_TPS_LIMIT_INTERVAL             = "60000"
	TPS_LIMIT_STRATEGY_KEY                 = "tps.limit.strategy"
	EXECUTE_LIMIT_KEY                      = "execute.limit"
	DEFAULT_EXECUTE_LIMIT                  = "-1"
	EXECUTE_REJECTED_EXECUTION_HANDLER_KEY = "execute.limit.rejected.handler"
	PROVIDER_SHUTDOWN_FILTER               = "pshutdown"
	CONSUMER_SHUTDOWN_FILTER               = "cshutdown"
	SERIALIZATION_KEY                      = "serialization"
	PID_KEY                                = "pid"
	SYNC_REPORT_KEY                        = "sync.report"
	RETRY_PERIOD_KEY                       = "retry.period"
	RETRY_TIMES_KEY                        = "retry.times"
	CYCLE_REPORT_KEY                       = "cycle.report"
	DEFAULT_BLACK_LIST_RECOVER_BLOCK       = 16
)

const (
	APPLICATION_KEY          = "application"
	ORGANIZATION_KEY         = "organization"
	NAME_KEY                 = "name"
	MODULE_KEY               = "module"
	APP_VERSION_KEY          = "app.version"
	OWNER_KEY                = "owner"
	ENVIRONMENT_KEY          = "environment"
	METHOD_KEY               = "method"
	METHOD_KEYS              = "methods"
	RULE_KEY                 = "rule"
	RUNTIME_KEY              = "runtime"
	BACKUP_KEY               = "backup"
	ROUTERS_CATEGORY         = "routers"
	ROUTE_PROTOCOL           = "route"
	CONDITION_ROUTE_PROTOCOL = "condition"
	TAG_ROUTE_PROTOCOL       = "tag"
	PROVIDERS_CATEGORY       = "providers"
	ROUTER_KEY               = "router"
	EXPORT_KEY               = "export"
)

const (
	CONFIG_NAMESPACE_KEY  = "config.namespace"
	CONFIG_GROUP_KEY      = "config.group"
	CONFIG_APP_ID_KEY     = "config.appId"
	CONFIG_CLUSTER_KEY    = "config.cluster"
	CONFIG_CHECK_KEY      = "config.check"
	CONFIG_TIMEOUT_KET    = "config.timeout"
	CONFIG_LOG_DIR_KEY    = "config.logDir"
	CONFIG_VERSION_KEY    = "configVersion"
	COMPATIBLE_CONFIG_KEY = "compatible_config"
)

// Use for router module
const (
	// ConditionRouterName Specify file condition router name
	ConditionRouterName = "condition"
	// ConditionAppRouterName Specify listenable application router name
	ConditionAppRouterName = "app"
	// ListenableRouterName Specify listenable router name
	ListenableRouterName = "listenable"
	// HealthCheckRouterName Specify the name of HealthCheckRouter
	HealthCheckRouterName = "health_check"
	// LocalPriorityRouterName Specify the name of LocalPriorityRouter
	LocalPriorityRouterName = "local_priority"
	// ConnCheckRouterName Specify the name of ConnCheckRouter
	ConnCheckRouterName = "conn_check"
	// TagRouterName Specify the name of TagRouter
	TagRouterName = "tag"
	// TagRouterRuleSuffix Specify tag router suffix
	TagRouterRuleSuffix  = ".tag-router"
	RemoteApplicationKey = "remote.application"
	// ConditionRouterRuleSuffix Specify condition router suffix
	ConditionRouterRuleSuffix = ".condition-router"

	// Force Force key in router module
	RouterForce = "force"
	// Enabled Enabled key in router module
	RouterEnabled = "enabled"
	// Priority Priority key in router module
	RouterPriority = "priority"
	// RouterScope Scope key in router module
	RouterScope = "scope"
	// RouterApplicationScope Scope key in router module
	RouterApplicationScope = "application"
	// RouterServiceScope Scope key in router module
	RouterServiceScope = "service"
	// RouterRuleKey defines the key of the router, service's/application's name
	RouterRuleKey = "key"
	// ForceUseTag is the tag in attachment
	ForceUseTag = "dubbo.force.tag"
	Tagkey      = "dubbo.tag"
	// HEALTH_ROUTE_ENABLED_KEY defines if use health router
	HEALTH_ROUTE_ENABLED_KEY = "health.route.enabled"
	// AttachmentKey in context in invoker
	AttachmentKey = DubboCtxKey("attachment")
)
const (
	ASYNC_KEY = "async" // it's value should be "true" or "false" of string type
)
const (
	SERVICE_FILTER_KEY   = "service.filter"
	REFERENCE_FILTER_KEY = "reference.filter"
)

const (
	// name of consumer sign filter
	CONSUMER_SIGN_FILTER = "sign"
	// name of consumer sign filter
	PROVIDER_AUTH_FILTER = "auth"
	// name of service filter
	SERVICE_AUTH_KEY = "auth"
	// key of authenticator
	AUTHENTICATOR_KEY = "authenticator"
	// name of default authenticator
	DEFAULT_AUTHENTICATOR = "accesskeys"
	// name of default url storage
	DEFAULT_ACCESS_KEY_STORAGE = "urlstorage"
	// key of storage
	ACCESS_KEY_STORAGE_KEY = "accessKey.storage"
	// key of request timestamp
	REQUEST_TIMESTAMP_KEY = "timestamp"
	// key of request signature
	REQUEST_SIGNATURE_KEY = "signature"
	// AK key
	AK_KEY = "ak"
	// signature format
	SIGNATURE_STRING_FORMAT = "%s#%s#%s#%s"
	// key whether enable signature
	PARAMETER_SIGNATURE_ENABLE_KEY = "param.sign"
	// consumer
	CONSUMER = "consumer"
	// key of access key id
	ACCESS_KEY_ID_KEY = ".accessKeyId"
	// key of secret access key
	SECRET_ACCESS_KEY_KEY = ".secretAccessKey"
)
const (
	NACOS_KEY                    = "nacos"
	NACOS_DEFAULT_ROLETYPE       = 3
	NACOS_CACHE_DIR_KEY          = "cacheDir"
	NACOS_LOG_DIR_KEY            = "logDir"
	NACOS_ENDPOINT               = "endpoint"
	NACOS_SERVICE_NAME_SEPARATOR = ":"
	NACOS_CATEGORY_KEY           = "category"
	NACOS_PROTOCOL_KEY           = "protocol"
	NACOS_PATH_KEY               = "path"
	NACOS_NAMESPACE_ID           = "namespaceId"
	NACOS_PASSWORD               = "password"
	NACOS_USERNAME               = "username"
	NACOS_NOT_LOAD_LOCAL_CACHE   = "nacos.not.load.cache"
)
const (
	TRACING_REMOTE_SPAN_CTX = DubboCtxKey("tracing.remote.span.ctx")
)
