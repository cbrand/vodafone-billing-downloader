package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type BearerToken interface {
	AuthenticateAPI(request *http.Request) *http.Request
}

var (
	ErrJsonRequestFailed = errors.New("JSON request failed")
)

func GetJson(url string, bearerToken BearerToken, target interface{}) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request = setJsonRequestHeaders(request)
	request = bearerToken.AuthenticateAPI(request)
	return jsonFromRequest(request, target)
}

func setJsonRequestHeaders(request *http.Request) *http.Request {
	request.Header.Set("Content-Type", "application/json")
	return request
}

func jsonFromRequest(request *http.Request, target interface{}) error {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		requestBytes, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(requestBytes))
		fmt.Println("")

		responseBytes, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(responseBytes))
		return ErrJsonRequestFailed
	}

	return getJsonFromResponse(response, target)
}

func getJsonFromResponse(response *http.Response, target interface{}) error {
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
