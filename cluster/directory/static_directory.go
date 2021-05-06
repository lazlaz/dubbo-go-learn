package directory

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/protocol"
	"go.uber.org/atomic"
)

type staticDirectory struct {
	BaseDirectory
	invokers []protocol.Invoker
}

func (staticDirectory) GetUrl() *common.URL {
	panic("implement me")
}

func (staticDirectory) IsAvailable() bool {
	panic("implement me")
}

func (staticDirectory) Destroy() {
	panic("implement me")
}

func (staticDirectory) List(invocation protocol.Invocation) []protocol.Invoker {
	panic("implement me")
}

// NewStaticDirectory Create a new staticDirectory with invokers
func NewStaticDirectory(invokers []protocol.Invoker) *staticDirectory {
	var url *common.URL

	if len(invokers) > 0 {
		url = invokers[0].GetUrl()
	}
	dir := &staticDirectory{
		BaseDirectory: NewBaseDirectory(url),
		invokers:      invokers,
	}

	return dir
}

// NewBaseDirectory Create BaseDirectory with URL
func NewBaseDirectory(url *common.URL) BaseDirectory {
	return BaseDirectory{
		url:       url,
		destroyed: atomic.NewBool(false),
	}
}
