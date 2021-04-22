package main

import (
	"go.uber.org/zap"
	"net/http"
)

var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger2()
	defer sugarLogger.Sync()
	simpleHttpGet2("www.google.com")
	simpleHttpGet2("http://www.baidu.com")
}

func InitLogger2() {
	logger, _ := zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet2(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
