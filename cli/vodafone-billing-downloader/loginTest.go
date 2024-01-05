package cli

import (
	"fmt"

	"github.com/cbrand/vodafone-billing-downloader/login"
	"github.com/urfave/cli/v2"
)

var LoginTestCommand = &cli.Command{
	Name:  "login-test",
	Usage: "Test your login credentials",
	Flags: []cli.Flag{
		FlagUsername,
		FlagPassword,
	},
	Action: func(c *cli.Context) error {
		_, err := login.Do(c.String(FlagUsername.Name), c.String(FlagPassword.Name))
		if err != nil {
			fmt.Println("Login failed")
			fmt.Println(err)
		} else {
			fmt.Println("Login successful")
		}
		return err
	},
}
