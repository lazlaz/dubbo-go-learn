package service

import "github.com/laz/dubbo-go/common"

type MetadataService interface {
	common.RPCService
	// ServiceName will get the service's name in meta service , which is application name
	ServiceName() (string, error)
	// ExportURL will store the exported url in metadata
	ExportURL(url *common.URL) (bool, error)
	// UnexportURL will delete the exported url in metadata
	UnexportURL(url *common.URL) error
	// SubscribeURL will store the subscribed url in metadata
	SubscribeURL(url *common.URL) (bool, error)
	// UnsubscribeURL will delete the subscribed url in metadata
	UnsubscribeURL(url *common.URL) error
	// PublishServiceDefinition will generate the target url's code info
	PublishServiceDefinition(url *common.URL) error
	// GetExportedURLs will get the target exported url in metadata
	// the url should be unique
	// due to dubbo-go only support return array []interface{} in RPCService, so we should declare the return type as []interface{}
	// actually, it's []String
	GetExportedURLs(serviceInterface string, group string, version string, protocol string) ([]interface{}, error)

	MethodMapper() map[string]string

	// GetExportedURLs will get the target subscribed url in metadata
	// the url should be unique
	GetSubscribedURLs() ([]*common.URL, error)
	// GetServiceDefinition will get the target service info store in metadata
	GetServiceDefinition(interfaceName string, group string, version string) (string, error)
	// GetServiceDefinition will get the target service info store in metadata by service key
	GetServiceDefinitionByServiceKey(serviceKey string) (string, error)
	// RefreshMetadata will refresh the metadata
	RefreshMetadata(exportedRevision string, subscribedRevision string) (bool, error)
	// Version will return the metadata service version
	Version() (string, error)
}
