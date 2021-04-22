package main

import (
	"github.com/laz/dubbo-go/config"
	"github.com/laz/dubbo-go/dubbo-test/pkg"
)
import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

func main() {
	hessian.RegisterPOJO(&pkg.User{})
	config.Load()

}
