package cli

import "github.com/urfave/cli/v2"

var (
	FlagUsername = &cli.StringFlag{
		Name:    "username",
		Usage:   "Your Vodafone username",
		EnvVars: []string{"VODAFONE_USERNAME"},
	}
	FlagPassword = &cli.StringFlag{
		Name:    "password",
		Usage:   "Your Vodafone password",
		EnvVars: []string{"VODAFONE_PASSWORD"},
	}
	FlagDirectory = &cli.PathFlag{
		Name:    "directory",
		Usage:   "The directory to store the invoices in",
		EnvVars: []string{"VODAFONE_DUMP_DIRECTORY"},
	}
	FlagPaperlessURL = &cli.StringFlag{
		Name:    "paperless-url",
		Usage:   "The URL of the paperless instance to upload the invoices to",
		EnvVars: []string{"PAPERLESS_URL"},
	}
	FlagPaperlessToken = &cli.StringFlag{
		Name:    "paperless-token",
		Usage:   "The token to use for authentication against paperless",
		EnvVars: []string{"PAPERLESS_TOKEN"},
	}
	FlagPaperlessUsername = &cli.StringFlag{
		Name:    "paperless-username",
		Usage:   "The username to use for authentication against paperless. Will be ignored if token is set.",
		EnvVars: []string{"PAPERLESS_USERNAME"},
	}
	FlagPaperlessPassword = &cli.StringFlag{
		Name:    "paperless-password",
		Usage:   "The password to use for authentication against paperless. Will be ignored if token is set.",
		EnvVars: []string{"PAPERLESS_PASSWORD"},
	}
	FlagPaperlessCorrespondent = &cli.StringFlag{
		Name:    "paperless-correspondent",
		Usage:   "The ID of correspondent to use for the invoices",
		EnvVars: []string{"PAPERLESS_CORRESPONDENT"},
	}
	FlagPaperlessDocumentType = &cli.StringFlag{
		Name:    "paperless-document-type",
		Usage:   "The ID of the document type to use for the invoices",
		EnvVars: []string{"PAPERLESS_DOCUMENT_TYPE"},
	}
)
