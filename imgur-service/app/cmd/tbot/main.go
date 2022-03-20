package main

import (
	"flag"

	"github.com/serj1c/tbot/internal"
	"github.com/serj1c/tbot/internal/config"
	"github.com/serj1c/tbot/pkg/logger"
)

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "config", "configs/dev.yml", "config file path")
}

func main() {
	flag.Parse()

	cfg := config.GetConfig(cfgPath)

	logger.Init(cfg.AppConfig.LogLevel)
	logger := logger.GetLogger()

	app, err := internal.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	app.Run()
}
