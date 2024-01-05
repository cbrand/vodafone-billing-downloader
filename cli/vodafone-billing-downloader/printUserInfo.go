package cli

import (
	"fmt"

	"github.com/cbrand/vodafone-billing-downloader/login"
	"github.com/urfave/cli/v2"
)

var UserInfoCommand = &cli.Command{
	Name:  "user-info",
	Usage: "List information about the logged in users.",
	Flags: []cli.Flag{
		FlagUsername,
		FlagPassword,
	},
	Action: func(c *cli.Context) error {
		loginData, err := login.Do(c.String(FlagUsername.Name), c.String(FlagPassword.Name))
		if err != nil {
			fmt.Println("Login failed")
			return err
		}

		userInfo, err := login.GetUserInfo(loginData)
		if err != nil {
			fmt.Println("Failed to retrieve user info")
			return err
		}

		fmt.Println(userInfo.HumanReadableString())

		return err
	},
}
