package main

import (
	"os"

	commands "github.com/cbrand/vodafone-billing-downloader/cli/vodafone-billing-downloader"
)

const Version = "0.0.1"

func main() {
	cli := commands.GetCLI(Version)
	err := cli.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
