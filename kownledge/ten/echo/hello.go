package echo

import (
	"github.com/apache/dubbo-getty"
	"github.com/laz/dubbo-go/kownledge/ten/echo/handler"
)

var (
	Sessions []getty.Session
)

func ClientRequest() {
	for _, session := range Sessions {
		ss := session
		go func() {
			echoTimes := 10
			for i := 0; i < echoTimes; i++ {
				_, _, err := ss.WritePkg("hello", handler.WritePkgTimeout)
				if err != nil {
					handler.Log.Infof("session.WritePkg(session{%s}, error{%v}", ss.Stat(), err)
					ss.Close()
				}
			}
			handler.Log.Infof("after loop %d times", echoTimes)
		}()
	}

}
