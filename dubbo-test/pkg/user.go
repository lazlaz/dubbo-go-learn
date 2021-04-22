package pkg

import (
	"context"
	gxlog "github.com/dubbogo/gost/log"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/config"

	"time"
)

func init() {
	config.SetProviderService(new(UserProvider))

}

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, req []interface{}) (*User, error) {
	gxlog.CInfo("req:%#v", req)

	t := time.Now()
	attachment := ctx.Value(constant.AttachmentKey).(map[string]interface{})
	if v, ok := attachment["timestamp"]; ok {
		gxlog.CInfo("attachment: %v", v)
		t = v.(time.Time).Add(-1 * 365 * 24 * time.Hour)
	}

	rsp := User{"A001", "Alex Stocks", 18, t}
	gxlog.CInfo("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

func (u User) JavaClassName() string {
	return "org.apache.dubbo.User"
}
