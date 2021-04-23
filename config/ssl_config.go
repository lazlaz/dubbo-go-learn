package config

import getty "github.com/apache/dubbo-getty"

var (
	serverTlsConfigBuilder getty.TlsConfigBuilder
	clientTlsConfigBuilder getty.TlsConfigBuilder
)

func GetServerTlsConfigBuilder() getty.TlsConfigBuilder {
	return serverTlsConfigBuilder
}
