package invoice

import (
	"testing"

	"github.com/cbrand/vodafone-billing-downloader/login"
	"github.com/stretchr/testify/assert"
)

func TestGetAndDownload(t *testing.T) {
	credentials, err := login.FromEnvironmentVariables()
	assert.Nil(t, err)
	assert.NotNil(t, credentials)
	userInfo, err := login.GetUserInfo(credentials)
	assert.Nil(t, err)
	assert.NotNil(t, userInfo)
	invoices, err := List(userInfo, credentials)
	assert.Nil(t, err)
	assert.NotNil(t, invoices)
	for _, invoiceList := range invoices.Data {
		for _, invoice := range invoiceList.Invoices {
			for _, document := range invoice.Documents {
				data, err := document.Download(credentials)
				assert.Nil(t, err)
				payload, err := data.Bytes()
				assert.Nil(t, err)
				assert.NotNil(t, payload)
			}
		}
	}
}
