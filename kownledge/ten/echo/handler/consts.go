package handler

import getty "github.com/apache/dubbo-getty"

const (
	CronPeriod      = 20e9
	WritePkgTimeout = 1e8
)

var (
	Log = getty.GetLogger()
)
