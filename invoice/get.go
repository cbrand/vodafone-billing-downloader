package invoice

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/cbrand/vodafone-billing-downloader/fetcher"
)

const (
	INVOICE_URL_TEMPLATE          = "https://api.vodafone.de/meinvodafone/v2/customer/urn:vf-de:cable:can:%s/invoice"
	INVOICE_DOCUMENT_URL_TEMPLATE = "https://api.vodafone.de/meinvodafone/v2/customer/%s/invoiceDocument/%s"
)

type ContractID string

func invoiceURLFor(contractID ContractID) string {
	return fmt.Sprintf(INVOICE_URL_TEMPLATE, contractID)
}

type ContractIDFetcher interface {
	GetAllContractIDs() []string
}

type InvoiceOverview struct {
	Data map[ContractID]*InvoiceList
}

func (invoiceOverview *InvoiceOverview) GetNumInvoices() int {
	numInvoices := 0
	for _, invoiceList := range invoiceOverview.Data {
		numInvoices += invoiceList.GetNumInvoices()
	}
	return numInvoices
}

func (invoiceOverview *InvoiceOverview) GetNumDocuments() int {
	numDocuments := 0
	for _, invoiceList := range invoiceOverview.Data {
		numDocuments += invoiceList.GetNumDocuments()
	}
	return numDocuments
}

func List(contractIDFetcher ContractIDFetcher, bearerToken fetcher.BearerToken) (*InvoiceOverview, error) {
	invoices := make(map[ContractID]*InvoiceList)
	for _, contractID := range contractIDFetcher.GetAllContractIDs() {
		invoiceList, err := ListFor(ContractID(contractID), bearerToken)
		if err != nil {
			return nil, err
		}
		invoices[ContractID(contractID)] = invoiceList
	}

	return &InvoiceOverview{Data: invoices}, nil
}

func ListFor(contractID ContractID, bearerToken fetcher.BearerToken) (*InvoiceList, error) {
	invoiceList := &InvoiceList{}
	err := fetcher.GetJson(invoiceURLFor(contractID), bearerToken, invoiceList)
	if err != nil {
		return nil, err
	}
	invoiceList.PropagateCustomerID()
	return invoiceList, nil
}

type InvoiceList struct {
	CustomerID string     `json:"customerId"`
	Invoices   []*Invoice `json:"invoices"`
}

func (invoiceList *InvoiceList) GetNumInvoices() int {
	return len(invoiceList.Invoices)
}

func (invoiceList *InvoiceList) GetNumDocuments() int {
	numDocuments := 0
	for _, invoice := range invoiceList.Invoices {
		numDocuments += invoice.GetNumDocuments()
	}
	return numDocuments
}

func (invoiceList *InvoiceList) PropagateCustomerID() {
	for _, invoice := range invoiceList.Invoices {
		invoice.SetCustomerID(invoiceList.CustomerID)
	}
}

type Invoice struct {
	Number     string             `json:"number"`
	Date       string             `json:"date"`
	Amount     float64            `json:"amount"` // Vodafone API has this defined as a float in json. This is a bad idea.
	DueDate    string             `json:"dueDate"`
	From       string             `json:"from"`
	About      string             `json:"about"`
	Documents  []*InvoiceDocument `json:"documents"`
	CustomerID string             `json:"-"`
}

func (invoice *Invoice) GetNumDocuments() int {
	return len(invoice.Documents)
}

func (invoice *Invoice) SetCustomerID(customerID string) {
	invoice.CustomerID = customerID
	for _, document := range invoice.Documents {
		document.SetCustomerID(customerID)
	}
}

type InvoiceDocument struct {
	DocumentID string `json:"documentId"`
	Category   string `json:"category"`
	Icon       string `json:"icon"`
	SubType    string `json:"subType"`
	CustomerID string `json:"-"`
}

func (invoiceDocument *InvoiceDocument) SetCustomerID(customerID string) {
	invoiceDocument.CustomerID = customerID
}

func (invoiceDocument *InvoiceDocument) DownloadURL() string {
	return fmt.Sprintf(INVOICE_DOCUMENT_URL_TEMPLATE, invoiceDocument.CustomerID, invoiceDocument.DocumentID)
}

func (invoiceDocument *InvoiceDocument) Download(bearerToken fetcher.BearerToken) (*DocumentData, error) {
	documentData := &DocumentData{}
	invoiceDownloadURL := invoiceDocument.DownloadURL()
	err := fetcher.GetJson(invoiceDownloadURL, bearerToken, documentData)
	if err != nil {
		return nil, err
	}
	return documentData, nil
}

type DocumentData struct {
	CustomerID string `json:"customerId"`
	DocumentID string `json:"documentId"`
	MimeType   string `json:"mime"`
	Data       string `json:"data"`
}

func (documentData *DocumentData) Bytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(documentData.Data)
}

// Checksum returns an md5checksum the same as paperless would calculate for
// cross referencing.
func (documentData *DocumentData) Checksum() string {
	data, _ := documentData.Bytes()
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
