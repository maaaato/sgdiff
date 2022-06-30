package cli

import (
	"github.com/maaaato/sgdiff/pkg/sgdiff"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	app := cli.NewApp()

	app.Name = "aaaa"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		sgdiff.NewCommand(),
	}
	return app
}
