package main

import (
	"github.com/shiranr/linkcheck/models"
	"github.com/shiranr/linkcheck/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// TODO mock external server behavior
func main() {
	start := time.Now()
	log.SetLevel(log.InfoLevel)
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := basepath + "/configuration/linkcheck.json"
	var app = &cli.App{
		Name:  "linkcheck",
		Usage: "A linter in Golang to verify Markdown links.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       configPath,
				Usage:       "configuration file",
				Destination: &configPath,
			},
		},
		Version: "1.0.0",
		Action: func(ctx *cli.Context) error {
			configPath = ctx.String("config")
			utils.LoadConfiguration(configPath)
			var readmeFiles []string
			log.Info("Context first argument " + ctx.Args().First())
			if ctx.Args().Present() {
				readmeFiles = utils.ExtractMarkdownFilesFromList(ctx.Args().Slice())
			} else {
				readmeFiles = utils.ExtractMarkdownFiles()
			}
			return models.GetFilesProcessorInstance().Process(readmeFiles)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	end := time.Now()
	log.Info("Time elapsed: " + end.Sub(start).String())
}
