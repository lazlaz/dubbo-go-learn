package protocol

import (
	"github.com/laz/dubbo-go/common"
	"github.com/laz/dubbo-go/common/logger"
	"sync"
)

// Protocol
// Extension - protocol
type Protocol interface {
	// Export service for remote invocation
	Export(invoker Invoker) Exporter
	// Refer a remote service
	Refer(url *common.URL) Invoker
	// Destroy will destroy all invoker and exporter, so it only is called once.
	Destroy()
}

// BaseProtocol is default protocol implement.
type BaseProtocol struct {
	exporterMap *sync.Map
	invokers    []Invoker
}

// Exporter
// wrapping invoker
type Exporter interface {
	// GetInvoker gets invoker.
	GetInvoker() Invoker
	// Unexport exported service.
	Unexport()
}

// NewBaseProtocol creates a new BaseProtocol
func NewBaseProtocol() BaseProtocol {
	return BaseProtocol{
		exporterMap: new(sync.Map),
	}
}

// ExporterMap gets exporter map.
func (bp *BaseProtocol) ExporterMap() *sync.Map {
	return bp.exporterMap
}

// SetExporterMap set @exporter with @key to local memory.
func (bp *BaseProtocol) SetExporterMap(key string, exporter Exporter) {
	bp.exporterMap.Store(key, exporter)
}

// BaseExporter is default exporter implement.
type BaseExporter struct {
	key         string
	invoker     Invoker
	exporterMap *sync.Map
}

// GetInvoker gets invoker
func (de *BaseExporter) GetInvoker() Invoker {
	return de.invoker

}

// Unexport exported service.
func (de *BaseExporter) Unexport() {
	logger.Infof("Exporter unexport.")
	de.invoker.Destroy()
	de.exporterMap.Delete(de.key)
}

// NewBaseExporter creates a new BaseExporter
func NewBaseExporter(key string, invoker Invoker, exporterMap *sync.Map) *BaseExporter {
	return &BaseExporter{
		key:         key,
		invoker:     invoker,
		exporterMap: exporterMap,
	}
}
