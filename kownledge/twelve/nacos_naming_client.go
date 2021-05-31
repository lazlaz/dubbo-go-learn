package main

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//create clientConfig
	clientConfig := constant.ClientConfig{}

	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "127.0.0.1",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}
	// Another way of create naming client for service discovery (recommend)
	namingClient, _ := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	_ = namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: "seata-server",
		GroupName:   "SEATA_GROUP", // default value is DEFAULT_GROUP
		Clusters:    []string{},    // default value is DEFAULT

		SubscribeCallback: func(services []model.SubscribeService, err error) {
			data, err := json.Marshal(services)
			log.Printf("servie info change：" + string(data))
		},
	})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig

}
