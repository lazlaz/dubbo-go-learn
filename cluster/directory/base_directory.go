package directory

import (
	"github.com/laz/dubbo-go/common"
	"go.uber.org/atomic"
	"sync"
)

// BaseDirectory Abstract implementation of Directory: Invoker list returned from this Directory's list method have been filtered by Routers
type BaseDirectory struct {
	url       *common.URL
	destroyed *atomic.Bool
	// this mutex for change the properties in BaseDirectory, like routerChain , destroyed etc
	mutex sync.Mutex
}
