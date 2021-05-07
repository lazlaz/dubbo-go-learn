package protocol

import (
	"context"
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/logger"
)
import (
	perrors "github.com/pkg/errors"
	uatomic "go.uber.org/atomic"
)

// Invoker the service invocation interface for the consumer
//go:generate mockgen -source invoker.go -destination mock/mock_invoker.go  -self_package github.com/apache/dubbo-go/protocol/mock --package mock  Invoker
// Extension - Invoker
type Invoker interface {
	common.Node
	// Invoke the invocation and return result.
	Invoke(context.Context, Invocation) Result
}

var (
	// ErrClientClosed means client has clossed.
	ErrClientClosed = perrors.New("remoting client has closed")
	// ErrNoReply
	ErrNoReply = perrors.New("request need @response")
	// ErrDestroyedInvoker
	ErrDestroyedInvoker = perrors.New("request Destroyed invoker")
)

// BaseInvoker provides default invoker implement
type BaseInvoker struct {
	url       *common.URL
	available uatomic.Bool
	destroyed uatomic.Bool
}

// NewBaseInvoker creates a new BaseInvoker
func NewBaseInvoker(url *common.URL) *BaseInvoker {
	ivk := &BaseInvoker{
		url: url,
	}
	ivk.available.Store(true)
	ivk.destroyed.Store(false)

	return ivk
}

// GetUrl gets base invoker URL
func (bi *BaseInvoker) GetUrl() *common.URL {
	return bi.url
}

// IsAvailable gets available flag
func (bi *BaseInvoker) IsAvailable() bool {
	return bi.available.Load()
}

// IsDestroyed gets destroyed flag
func (bi *BaseInvoker) IsDestroyed() bool {
	return bi.destroyed.Load()
}

// Invoke provides default invoker implement
func (bi *BaseInvoker) Invoke(context context.Context, invocation Invocation) Result {
	return &RPCResult{}
}

// Destroy changes available and destroyed flag
func (bi *BaseInvoker) Destroy() {
	logger.Infof("Destroy invoker: %s", bi.GetUrl())
	bi.destroyed.Store(true)
	bi.available.Store(false)
}
