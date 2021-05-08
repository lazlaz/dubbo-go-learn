package impl

const (
	// Zero : byte zero
	Zero = byte(0x00)
)

/**
 * the dubbo protocol header length is 16 Bytes.
 * the first 2 Bytes is magic code '0xdabb'
 * the next 1 Byte is message flags, in which its 16-20 bit is serial id, 21 for event, 22 for two way, 23 for request/response flag
 * the next 1 Bytes is response state.
 * the next 8 Bytes is package DI.
 * the next 4 Bytes is package length.
 **/
const (
	// header length.
	HEADER_LENGTH = 16

	// magic header
	MAGIC      = uint16(0xdabb)
	MAGIC_HIGH = byte(0xda)
	MAGIC_LOW  = byte(0xbb)

	// message flag.
	FLAG_REQUEST = byte(0x80)
	FLAG_TWOWAY  = byte(0x40)
	FLAG_EVENT   = byte(0x20) // for heartbeat
	SERIAL_MASK  = 0x1f

	DUBBO_VERSION                          = "2.5.4"
	DUBBO_VERSION_KEY                      = "dubbo"
	DEFAULT_DUBBO_PROTOCOL_VERSION         = "2.0.2" // Dubbo RPC protocol version, for compatibility, it must not be between 2.0.10 ~ 2.6.2
	LOWEST_VERSION_FOR_RESPONSE_ATTACHMENT = 2000200
	DEFAULT_LEN                            = 8388608 // 8 * 1024 * 1024 default body max length
)

// Body map keys
var (
	DubboVersionKey = "dubboVersion"
	ArgsTypesKey    = "argsTypes"
	ArgsKey         = "args"
	ServiceKey      = "service"
	AttachmentsKey  = "attachments"
)

// ResponsePayload related consts
const (
	Response_OK                byte = 20
	Response_CLIENT_TIMEOUT    byte = 30
	Response_SERVER_TIMEOUT    byte = 31
	Response_BAD_REQUEST       byte = 40
	Response_BAD_RESPONSE      byte = 50
	Response_SERVICE_NOT_FOUND byte = 60
	Response_SERVICE_ERROR     byte = 70
	Response_SERVER_ERROR      byte = 80
	Response_CLIENT_ERROR      byte = 90

	// According to "java dubbo" There are two cases of response:
	// 		1. with attachments
	// 		2. no attachments
	RESPONSE_WITH_EXCEPTION                  int32 = 0
	RESPONSE_VALUE                           int32 = 1
	RESPONSE_NULL_VALUE                      int32 = 2
	RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS int32 = 3
	RESPONSE_VALUE_WITH_ATTACHMENTS          int32 = 4
	RESPONSE_NULL_VALUE_WITH_ATTACHMENTS     int32 = 5
)

// Dubbo request response related consts
var (
	DubboRequestHeaderBytesTwoWay = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST | FLAG_TWOWAY}
	DubboRequestHeaderBytes       = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST}
	DubboResponseHeaderBytes      = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, Zero, Response_OK}
	DubboRequestHeartbeatHeader   = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST | FLAG_TWOWAY | FLAG_EVENT}
	DubboResponseHeartbeatHeader  = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_EVENT}
)
