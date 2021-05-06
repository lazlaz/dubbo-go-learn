package main

import (
	"context"
	hessian "github.com/apache/dubbo-go-hessian2"
	gxlog "github.com/dubbogo/gost/log"
	"github.com/laz/dubbo-go/config"
	pkg "github.com/laz/dubbo-go/dubbo-test/pkg/client"
	"os"
	"time"
)

var userProvider = new(pkg.UserProvider)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

func main() {
	config.Load()
	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &pkg.User{}
	err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %v\n", user)
}