package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/koolay/db-replay/proxy"
	"github.com/sirupsen/logrus"
)

const defaultVersion = "5.7.28"

func newLogger(levelName string, hooks ...logrus.Hook) *logrus.Logger {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		level = logrus.ErrorLevel
	}

	logFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, fmt.Sprintf("%s:%v", filename, f.Line)
		},
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(logFormatter)
	logger.SetReportCaller(true)

	for _, hook := range hooks {
		logger.AddHook(hook)
	}
	return logger
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg := proxy.Config{
		Addr:     "0.0.0.0:3309",
		User:     "root",
		Password: "123",
		TargetConnection: proxy.Connection{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Password: "dev",
			Database: "myapp",
		},
	}

	logger := newLogger("DEBUG")
	ser := proxy.NewProxyServer(defaultVersion, cfg)
	handler, err := proxy.NewMysqlHandler(cfg.TargetConnection, logger)
	if err != nil {
		logger.WithError(err).Error("failed to new handler")
		return
	}

	proxier, err := proxy.NewProxy(cfg, logger, ser, handler)
	if err != nil {
		logger.WithError(err).Error("failed to new proxy")
		return
	}

	go func() {
		log.Println("start db proxy")
		if err := proxier.Start(); err != nil {
			logger.WithError(err).Error("failed to start proxy")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := proxier.Shutdown(); err != nil {
		log.Println(err)
	}
	log.Println("shutdown")
}
