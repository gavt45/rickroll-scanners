package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/gavt45/rickroll-scanners/pkg/app"
	"github.com/gavt45/rickroll-scanners/pkg/config"
)

var (
	BuildDate = "unknown"
	Version   = "unknown"
	GitCommit = "unknown"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to config")
	flag.Parse()

	cfg, err := config.ReadConfig(configPath)
	if errors.Is(err, os.ErrNotExist) {
		cfg = config.DefaultConfig
	} else if err != nil {
		log.Println("error reading config: ", err.Error())

		return
	}

	app := app.New(cfg)

	log.Println(
		"Rickrolls standalone starting!",
		"Version:", Version,
		"Build date:", BuildDate,
		"Commit:", GitCommit,
	)

	if err := app.Start(context.Background()); err != nil {
		log.Println("error running app: ", err.Error())

		return
	}
}
