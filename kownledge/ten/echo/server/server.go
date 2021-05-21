package main

import (
	"flag"
	"github.com/apache/dubbo-getty/demo/util"
	tcp "github.com/laz/dubbo-go/kownledge/ten/echo"
)

import (
	"github.com/apache/dubbo-getty"
	gxsync "github.com/dubbogo/gost/sync"
)

var (
	taskPoolMode = flag.Bool("taskPool", false, "task pool mode")
	taskPoolSize = flag.Int("task_pool_size", 2000, "task poll size")
	pprofPort    = flag.Int("pprof_port", 65432, "pprof http port")
)

var taskPool gxsync.GenericTaskPool

func main() {
	flag.Parse()

	//util.SetLimit()
	//
	//util.Profiling(*pprofPort)

	options := []getty.ServerOption{getty.WithLocalAddress(":8090")}

	if *taskPoolMode {
		taskPool = gxsync.NewTaskPoolSimple(*taskPoolSize)
		options = append(options, getty.WithServerTaskPool(taskPool))
	}

	server := getty.NewTCPServer(options...)

	go server.RunEventLoop(NewHelloServerSession)

	util.WaitCloseSignals(server)
}

func NewHelloServerSession(session getty.Session) (err error) {
	err = tcp.InitialSession(session)
	if err != nil {
		return
	}
	return
}
