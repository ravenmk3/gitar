package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	AppName = "gitar"
)

func RunCliApp() error {
	app := NewCliApp()
	return app.Run(os.Args)
}

func NewCliApp() *cli.App {
	app := &cli.App{
		Name:        AppName,
		Usage:       AppName,
		Description: "Git Archive & Repo Tool",
		Commands: []*cli.Command{
			NewDownloadCommand(),
		},
	}
	return app
}

func NewDownloadCommand() *cli.Command {
	return &cli.Command{
		Name:    "dl",
		Aliases: []string{"download"},
		Usage:   "Download git archive",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Required: false, Value: false},
			&cli.BoolFlag{Name: "mail", Aliases: []string{"m"}, Required: false, Value: false},
		},
		Action: func(ctx *cli.Context) error {
			debug := ctx.Bool("debug")
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			url := ctx.Args().First()
			return DownloadArchive(url, ctx.Bool("mail"))
		},
	}
}
