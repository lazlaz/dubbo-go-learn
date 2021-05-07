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

// GetUrl Get URL
func (dir *BaseDirectory) GetUrl() *common.URL {
	return dir.url
}

// IsAvailable Once directory init finish, it will change to true
func (dir *BaseDirectory) IsAvailable() bool {
	return !dir.destroyed.Load()
}

// Destroy Destroy
func (dir *BaseDirectory) Destroy(doDestroy func()) {
	if dir.destroyed.CAS(false, true) {
		dir.mutex.Lock()
		doDestroy()
		dir.mutex.Unlock()
	}
}

// GetDirectoryUrl Get URL instance
func (dir *BaseDirectory) GetDirectoryUrl() *common.URL {
	return dir.url
}
