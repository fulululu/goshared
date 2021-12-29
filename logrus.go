package goshared

import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogrusRawConfig struct {
	LogLevel  logrus.Level `envconfig:"-"`
	LogFile   string       `envconfig:"LOGRUS_LOG_FILE,optional"`
	LogCaller bool         `envconfig:"LOGRUS_LOG_CALLER,optional"`
}

// InitLogrus ...
func InitLogrus(cfg LogrusRawConfig) {
	var writters []io.Writer
	writters = append(writters, os.Stdout)
	if len(cfg.LogFile) != 0 {
		f, err := os.OpenFile(cfg.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		writters = append(writters, f)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		DisableQuote:  true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			filename = filename + ":" + strconv.Itoa(f.Line)
			return funcname, filename
		},
	})
	logrus.SetOutput(io.MultiWriter(writters...))
	logrus.SetLevel(cfg.LogLevel)
	logrus.SetReportCaller(cfg.LogCaller)
}

func NewLogrusLogger(cfg LogrusRawConfig) *logrus.Logger {
	var writters []io.Writer

	writters = append(writters, os.Stdout)
	if len(cfg.LogFile) != 0 {
		f, err := os.OpenFile(cfg.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		writters = append(writters, f)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		DisableQuote:  true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	})
	logger.SetOutput(io.MultiWriter(writters...))
	logger.SetLevel(cfg.LogLevel)
	logger.SetReportCaller(cfg.LogCaller)

	return logger
}
