package main

import (
	"fmt"
	"os"

	"github.com/maaaato/sgdiff/pkg/cli"
)

const (
	exitCodeOk int = iota
	exitCodeError
	exitCodeFileError
)

func main() {
	app := cli.NewApp()
	err := app.Run(os.Args)
	code := handleExit(err)
	os.Exit(code)
}

func handleExit(err error) int {
	if err == nil {
		return exitCodeOk
	}
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	return exitCodeError
}
