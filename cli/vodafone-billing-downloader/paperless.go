package cli

import (
	"bytes"
	"fmt"

	"github.com/cbrand/vodafone-billing-downloader/invoice"
	"github.com/cbrand/vodafone-billing-downloader/login"
	"github.com/cbrand/vodafone-billing-downloader/paperless"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

var DumpToPaperlessCommand = &cli.Command{
	Name:  "dump-to-paperless",
	Usage: "Dumps all invoices to paperless / paperless ngx.",
	Flags: []cli.Flag{
		FlagUsername,
		FlagPassword,
		FlagPaperlessURL,
		FlagPaperlessToken,
		FlagPaperlessUsername,
		FlagPaperlessPassword,
		FlagPaperlessCorrespondent,
		FlagPaperlessDocumentType,
	},
	Action: func(c *cli.Context) error {
		return DumpToPaperless(
			c.String(FlagUsername.Name),
			c.String(FlagPassword.Name),
			&PaperlessDumpConfig{
				Correspondent: c.String(FlagPaperlessCorrespondent.Name),
				DocumentType:  c.String(FlagPaperlessDocumentType.Name),
				Config: &paperless.Config{
					URL:      c.String(FlagPaperlessURL.Name),
					APIKey:   c.String(FlagPaperlessToken.Name),
					Username: c.String(FlagPaperlessUsername.Name),
					Password: c.String(FlagPaperlessPassword.Name),
				},
			},
		)
	},
}

type PaperlessDumpConfig struct {
	Config        *paperless.Config
	Correspondent string
	DocumentType  string
}

func DumpToPaperless(username string, password string, paperlessInfo *PaperlessDumpConfig) error {
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

	for customerID, invoiceList := range invoices.Data {
		for _, invoice := range invoiceList.Invoices {
			for _, document := range invoice.Documents {
				documentData, err := document.Download(loginData)
				if err != nil {
					return err
				}
				title := fmt.Sprintf("Vodafone Invoice %s %s", customerID, invoice.Date)
				fileName := fmt.Sprintf("vodafone-%s-%s-%s-%s.pdf", document.Category, customerID, invoice.Date, document.DocumentID)
				checkSumExists, err := paperless.ChecksumExists(paperlessInfo.Config, documentData.Checksum())
				if err != nil {
					fmt.Println("Could not check if checksum exists")
					fmt.Println(err)
					return err
				}
				if !checkSumExists {
					payload, err := documentData.Bytes()
					if err != nil {
						fmt.Println("Could not extract document data")
						return err
					}
					err = paperless.DumpTo(paperlessInfo.Config, &paperless.DocumentInformation{
						Title:         title,
						Created:       invoice.Date,
						Correspondent: paperlessInfo.Correspondent,
						DocumentType:  paperlessInfo.DocumentType,
						FileName:      fileName,
						Data:          bytes.NewBuffer(payload),
					})
					if err != nil {
						fmt.Println("Could not write document data")
						fmt.Println(err)
						return err
					}
				}
				progress.Add(1)
			}
		}
	}

	return nil
}
