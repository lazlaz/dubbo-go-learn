package interfaces

import "github.com/laz/dubbo-go/common"

// ConfigPostProcessor is an extension to give users a chance to customize configs against ReferenceConfig and
// ServiceConfig during deployment time.
type ConfigPostProcessor interface {
	// PostProcessReferenceConfig customizes ReferenceConfig's params.
	PostProcessReferenceConfig(*common.URL)

	// PostProcessServiceConfig customizes ServiceConfig's params.
	PostProcessServiceConfig(*common.URL)
}
