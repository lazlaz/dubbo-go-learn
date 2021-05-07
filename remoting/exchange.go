package remoting

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/logger"
	"go.uber.org/atomic"
	"sync"
	"time"
)

// this is request for transport layer
type Request struct {
	ID int64
	// protocol version
	Version string
	// serial ID (ignore)
	SerialID byte
	// Data
	Data   interface{}
	TwoWay bool
	Event  bool
}

// the client sends request to server, there is one pendingResponse at client side to wait the response from server
type PendingResponse struct {
	seq       int64
	Err       error
	start     time.Time
	ReadStart time.Time
	Callback  common.AsyncCallback
	response  *Response
	Reply     interface{}
	Done      chan struct{}
}

// this is response for transport layer
type Response struct {
	ID       int64
	Version  string
	SerialID byte
	Status   uint8
	Event    bool
	Error    error
	Result   interface{}
}

// NewResponse create to a new Response.
func NewResponse(id int64, version string) *Response {
	return &Response{
		ID:      id,
		Version: version,
	}
}

type SequenceType int64

func (response *Response) Handle() {
	pendingResponse := removePendingResponse(SequenceType(response.ID))
	if pendingResponse == nil {
		logger.Errorf("failed to get pending response context for response package %s", *response)
		return
	}

	pendingResponse.response = response

	if pendingResponse.Callback == nil {
		pendingResponse.Err = pendingResponse.response.Error
		close(pendingResponse.Done)
	} else {
		pendingResponse.Callback(pendingResponse.GetCallResponse())
	}
}

type Options struct {
	// connect timeout
	ConnectTimeout time.Duration
}

//AsyncCallbackResponse async response for dubbo
type AsyncCallbackResponse struct {
	common.CallbackResponse
	Opts      Options
	Cause     error
	Start     time.Time // invoke(call) start time == write start time
	ReadStart time.Time // read start time, write duration = ReadStart - Start
	Reply     interface{}
}

// GetCallResponse is used for callback of async.
// It is will return AsyncCallbackResponse.
func (r PendingResponse) GetCallResponse() common.CallbackResponse {
	return AsyncCallbackResponse{
		Cause:     r.Err,
		Start:     r.start,
		ReadStart: r.ReadStart,
		Reply:     r.response,
	}
}

var (
	// generate request ID for global use
	sequence atomic.Int64

	// store requestID and response
	pendingResponses = new(sync.Map)
)

// get and remove response
func removePendingResponse(seq SequenceType) *PendingResponse {
	if pendingResponses == nil {
		return nil
	}
	if presp, ok := pendingResponses.Load(seq); ok {
		pendingResponses.Delete(seq)
		return presp.(*PendingResponse)
	}
	return nil
}
