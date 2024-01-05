package paperless

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const API_DOCUMENT_SEARCH_PATH = "/api/documents/"

type checkSumResponse struct {
	Count int `json:"count"`
}

func ChecksumExists(config *Config, checksum string) (bool, error) {
	urlString, err := url.JoinPath(config.URL, API_DOCUMENT_SEARCH_PATH)
	if err != nil {
		return false, err
	}
	url, err := url.Parse(urlString)
	if err != nil {
		return false, err
	}
	query := url.Query()
	query.Set("query", fmt.Sprintf("checksum:%s", checksum))
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return false, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", config.GetAuthorizationHeader())

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	checkSumResponse := &checkSumResponse{}
	err = json.Unmarshal(responseBytes, checkSumResponse)
	if err != nil {
		return false, err
	}

	return checkSumResponse.Count > 0, nil
}
