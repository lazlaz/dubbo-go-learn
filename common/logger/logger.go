package logger

var (
	logger Logger
)

// Logger is the interface for Logger types
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
}

func init() {
	if logger != nil {
		return
	}
	//命令行解析
	//fs := flag.NewFlagSet("log",flag.ContinueOnError)
	//logConfFile := fs.String("logConf",os.Getenv(constant.APP_LOG_CONF_FILE),"default log config path")

}
