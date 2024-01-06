# Vodafone Billing Downloader

A script which dumps all cable internet bills from www.vodafone.de to a directory or uploads them to [paperless-ngx](https://docs.paperless-ngx.com/).

## Usage

The Vodafone Billing Downloader can be used with the following command line arguments:

_Generic Parameters_

- `--username`: Your Vodafone username. Can also be set with the `VODAFONE_USERNAME` environment variable.
- `--password`: Your Vodafone password. Can also be set with the `VODAFONE_PASSWORD` environment variable.

_Dump parameters_

- `--directory`: The directory to store the invoices in. Can also be set with the `VODAFONE_DUMP_DIRECTORY` environment variable.

_Paperless specific parameters_

- `--paperless-url`: The URL of the paperless instance to upload the invoices to. Can also be set with the `PAPERLESS_URL` environment variable.
- `--paperless-token`: The token to use for authentication against paperless. Can also be set with the `PAPERLESS_TOKEN` environment variable.
- `--paperless-username`: The username to use for authentication against paperless. Will be ignored if token is set. Can also be set with the `PAPERLESS_USERNAME` environment variable.
- `--paperless-password`: The password to use for authentication against paperless. Will be ignored if token is set. Can also be set with the `PAPERLESS_PASSWORD` environment variable.
- `--paperless-correspondent`: The ID of correspondent to use for the invoices. Can also be set with the `PAPERLESS_CORRESPONDENT` environment variable.
- `--paperless-document-type`: The ID of the document type to use for the invoices. Can also be set with the `PAPERLESS_DOCUMENT_TYPE` environment variable.

### Commands

- `login-test`: Test your login credentials
- `user-info`: List information about the logged in users.
- `dump`: Dump all invoices to the specified directory.
- `dump-to-paperless`: Dumps all invoices to paperless / paperless ngx.
- `help, h`: Shows a list of commands or help for one command

### Example

To dump all invoices to a specified directory, you can use the following command:

```bash
vodafone-billing-downloader --username your_username --password your_password --directory /path/to/directory dump
```

### Containers

Docker Containers are available no both [Docker Hub](https://hub.docker.com/r/cbrand/vodafone-billing-downloader) and [ghcr.io](https://ghcr.io/cbrand/vodafone-billing-downloader).
