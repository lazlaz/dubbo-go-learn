package main

import (
	_ "github.com/laz/dubbo-go/common/proxy/proxy_factory"
	"github.com/laz/dubbo-go/config"
	"github.com/laz/dubbo-go/dubbo-test/pkg"
	_ "github.com/laz/dubbo-go/protocol/dubbo"
	_ "github.com/laz/dubbo-go/registry/protocol"
)
import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

func main() {
	hessian.RegisterPOJO(&pkg.User{})
	config.Load()

}
