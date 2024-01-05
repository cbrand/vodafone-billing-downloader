package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/cbrand/vodafone-billing-downloader/invoice"
	"github.com/cbrand/vodafone-billing-downloader/login"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

var DumpInvoicesCommand = &cli.Command{
	Name:  "dump",
	Usage: "Dump all invoices to the specified directory.",
	Flags: []cli.Flag{
		FlagUsername,
		FlagPassword,
		FlagDirectory,
	},
	Action: func(c *cli.Context) error {
		return DumpInvoice(
			c.String(FlagUsername.Name),
			c.String(FlagPassword.Name),
			c.Path(FlagDirectory.Name),
		)
	},
}

func DumpInvoice(username string, password string, directory string) error {
	fmt.Println("Logging in")
	loginData, err := login.Do(username, password)
	if err != nil {
		fmt.Println("Login failed")
		return err
	}

	fmt.Println("Retrieving user info")
	userInfo, err := login.GetUserInfo(loginData)
	if err != nil {
		fmt.Println("Retrieving user info failed")
		return err
	}

	fmt.Println("Retrieving invoices")
	invoices, err := invoice.List(userInfo, loginData)
	if err != nil {
		fmt.Println("Retrieving invoices failed")
		return err
	}

	numDocuments := invoices.GetNumDocuments()
	fmt.Printf("Found %d documents\n", numDocuments)
	progress := progressbar.Default(int64(numDocuments))
	os.MkdirAll(directory, os.ModePerm)

	for customerID, invoiceList := range invoices.Data {
		for _, invoice := range invoiceList.Invoices {
			for _, document := range invoice.Documents {
				documentData, err := document.Download(loginData)
				if err != nil {
					return err
				}
				targetPath := path.Join(directory, fmt.Sprintf("%s-%s-%s.pdf", customerID, invoice.Date, document.DocumentID))

				payload, err := documentData.Bytes()
				if err != nil {
					fmt.Println("Could not extract document data")
					return err
				}
				err = os.WriteFile(targetPath, payload, os.ModePerm)
				if err != nil {
					fmt.Println("Could not write document data")
					return err
				}
				progress.Add(1)
			}
		}
	}

	return nil
}
