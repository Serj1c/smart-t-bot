package main

import (
	"flag"

	"github.com/serjic/yt-search/internal"
	"github.com/serjic/yt-search/internal/config"
	"github.com/serjic/yt-search/pkg/logger"
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

	app, err := internal.NewApp(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	app.Run()
}
