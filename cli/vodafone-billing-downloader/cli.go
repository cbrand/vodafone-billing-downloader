package cli

import (
	"github.com/urfave/cli/v2"
)

func GetCLI(Version string) *cli.App {
	app := cli.NewApp()
	app.Name = "vodafone-billing-downloader"
	app.Usage = "Download your Vodafone billing documents as Vodafone does not provide a way to download all of them at once or via an API."
	app.Version = Version
	app.Authors = []*cli.Author{
		{
			Name: "Christoph Brand",
		},
	}
	app.Commands = []*cli.Command{
		LoginTestCommand,
		UserInfoCommand,
		DumpInvoicesCommand,
		DumpToPaperlessCommand,
	}
	return app
}
