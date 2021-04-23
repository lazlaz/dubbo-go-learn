package extension

import (
	"fmt"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/metadata/service"
)
import (
	perrors "github.com/pkg/errors"
)

var (
	// there will be two types: local or remote
	metadataServiceInsMap = make(map[string]func() (service.MetadataService, error), 2)
	// remoteMetadataService
	remoteMetadataService service.MetadataService
)

// GetRemoteMetadataService will get a RemoteMetadataService instance
func GetRemoteMetadataService() (service.MetadataService, error) {
	if remoteMetadataService != nil {
		return remoteMetadataService, nil
	}
	if creator, ok := metadataServiceInsMap["remote"]; ok {
		var err error
		remoteMetadataService, err = creator()
		return remoteMetadataService, err
	}
	logger.Warn("could not find the metadata service creator for metadataType: remote")
	return nil, perrors.New(fmt.Sprintf("could not find the metadata service creator for metadataType: remote"))
}
