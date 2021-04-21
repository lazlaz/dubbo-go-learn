package logger

import (
	"flag"
	"github.com/laz/dubbo-go/common/constant"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"os"
	"path"
)
import (
	"github.com/apache/dubbo-getty"
	perrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	logger Logger
)

// nolint
type DubboLogger struct {
	Logger
	dynamicLevel zap.AtomicLevel
}

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
	fs := flag.NewFlagSet("log", flag.ContinueOnError)
	logConfFile := fs.String("logConf", os.Getenv(constant.APP_LOG_CONF_FILE), "default log config path")
	fs.Parse(os.Args[1:])
	for len(fs.Args()) != 0 {
		fs.Parse(fs.Args()[1:])
	}
	err := InitLog(*logConfFile)
	if err != nil {
		log.Printf("[InitLog] warn: %v", err)
	}
}
func InitLog(logConfFile string) error {
	if logConfFile == "" {
		InitLogger(nil)
		return perrors.New("log configure file name is nil")
	}
	if path.Ext(logConfFile) != ".yml" {
		InitLogger(nil)
		return perrors.Errorf("log configure file name{%s} suffix must be .yml", logConfFile)
	}

	confFileStream, err := ioutil.ReadFile(logConfFile)
	if err != nil {
		InitLogger(nil)
		return perrors.Errorf("ioutil.ReadFile(file:%s) = error:%v", logConfFile, err)
	}

	conf := &zap.Config{}
	err = yaml.Unmarshal(confFileStream, conf)
	if err != nil {
		InitLogger(nil)
		return perrors.Errorf("[Unmarshal]init logger error: %v", err)
	}

	InitLogger(conf)

	return nil
}

func InitLogger(conf *zap.Config) {
	var zapLoggerConfig zap.Config
	if conf == nil {
		zapLoggerEncoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		zapLoggerConfig = zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:      false,
			Encoding:         "console",
			EncoderConfig:    zapLoggerEncoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else {
		zapLoggerConfig = *conf
	}
	zapLogger, _ := zapLoggerConfig.Build(zap.AddCallerSkip(1))
	logger = &DubboLogger{Logger: zapLogger.Sugar(), dynamicLevel: zapLoggerConfig.Level}

	// set getty log
	getty.SetLogger(logger)
}

// SetLogger sets logger for dubbo and getty
func SetLogger(log Logger) {
	logger = log
	getty.SetLogger(logger)
}

// GetLogger gets the logger
func GetLogger() Logger {
	return logger
}
