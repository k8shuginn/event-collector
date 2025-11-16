package main

import (
	"os"
	"strconv"

	"github.com/k8shuginn/event-collector/cmd/collector/app"
	"github.com/k8shuginn/event-collector/cmd/collector/config"
	"github.com/k8shuginn/event-collector/pkg/logger"
	"github.com/k8shuginn/event-collector/pkg/pprof"
	"go.uber.org/zap"
)

const (
	AppName string = "event-collector"

	EnvLogLevel    string = "LOG_LEVEL"
	EnvLogSize     string = "LOG_SIZE"
	EnvLogAge      string = "LOG_AGE"
	EnvLogBack     string = "LOG_BACK"
	EnvLogCompress string = "LOG_COMPRESS"

	DefaultConfigPath string = "/etc/collector/config.yaml"
)

func init() {
	// set logger
	logLevel := os.Getenv(EnvLogLevel)
	logSize, _ := strconv.Atoi(os.Getenv(EnvLogSize))
	logAge, _ := strconv.Atoi(os.Getenv(EnvLogAge))
	logBack, _ := strconv.Atoi(os.Getenv(EnvLogBack))
	logCompress, _ := strconv.ParseBool(os.Getenv(EnvLogCompress))

	logger.CreateGlobalLogger(
		AppName,
		logger.WithLogLevel(logLevel),
		logger.WithLogMaxSize(logSize),
		logger.WithLogMaxAge(logAge),
		logger.WithLogMaxBackups(logBack),
		logger.WithLogCompress(logCompress),
	)

	// start pprof
	pprof.InitPprof()
}

func main() {
	cfg, err := config.LoadConfig(DefaultConfigPath)
	if err != nil {
		logger.Panic("failed to load config", zap.Error(err))
	}

	app, err := app.NewCollector(cfg)
	if err != nil {
		logger.Panic("failed to create application", zap.Error(err))
	}
	app.Run()
}
