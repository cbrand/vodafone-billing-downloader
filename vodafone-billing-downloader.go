package main

import (
	"os"

	commands "github.com/cbrand/vodafone-billing-downloader/cli/vodafone-billing-downloader"
)

var version string

func main() {
	cli := commands.GetCLI(version)
	err := cli.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
