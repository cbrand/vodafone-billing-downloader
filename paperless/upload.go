package paperless

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type DocumentInformation struct {
	Title         string
	Created       string
	Correspondent string
	DocumentType  string
	FileName      string
	Tags          []string
	Data          io.Reader
}

func (documentInformation *DocumentInformation) toMimeType(b *bytes.Buffer) *multipart.Writer {
	w := multipart.NewWriter(b)

	documentField, err := w.CreateFormFile("document", documentInformation.FileName)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(documentField, documentInformation.Data)
	if err != nil {
		panic(err)
	}

	if len(documentInformation.Title) > 0 {
		w.WriteField("title", documentInformation.Title)
	}
	if len(documentInformation.Created) > 0 {
		w.WriteField("created", documentInformation.Created)
	}
	if len(documentInformation.Correspondent) > 0 {
		w.WriteField("correspondent", documentInformation.Correspondent)
	}
	if len(documentInformation.DocumentType) > 0 {
		w.WriteField("document_type", documentInformation.DocumentType)
	}
	for _, tag := range documentInformation.Tags {
		w.WriteField("tags", tag)
	}

	return w
}

const API_DOCUMENT_UPLOAD_PATH = "/api/documents/post_document/"

func DumpTo(auth *Config, document *DocumentInformation) error {
	uploadURL, _ := url.JoinPath(auth.URL, API_DOCUMENT_UPLOAD_PATH)
	var b bytes.Buffer
	mimeTypeWriter := document.toMimeType(&b)
	request, err := http.NewRequest("POST", uploadURL, &b)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", mimeTypeWriter.FormDataContentType())
	request.Header.Add("Authorization", auth.GetAuthorizationHeader())

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		resp, _ := httputil.DumpResponse(response, true)
		fmt.Println(string(resp))
		return fmt.Errorf("paperless upload. unexpected status code: %d", response.StatusCode)
	}
}
