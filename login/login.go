package login

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	DEFAULT_ACCEPT             = "application/json, text/plain, */*"
	DEFAULT_ACCEPT_LANGUAGE    = "en-US,en;q=0.5"
	DEFAULT_LOGIN_URL          = "https://www.vodafone.de/mint/rest/v60/session/start"
	DEFAULT_TARGET_URL         = ""
	DEFAULT_USER_AGENT         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0"
	DEFAULT_REFERER            = "https://www.vodafone.de/meinvodafone/account/login"
	OIDC_AUTHORIZE_URL         = "https://www.vodafone.de/mint/oidc/authorize"
	OIDC_RESPONSE_TYPE         = "code"
	OIDC_CLIENT_ID             = "b0595a44-0726-11ec-9011-9457a55a403c"
	OIDC_SCOPES                = "openid profile webseal user-groups user-accounts validate-token update-email-username account"
	OIDC_REDIRECT_URI          = "https://www.vodafone.de/meinvodafone/services/"
	OIDC_CODE_CHALLENGE_METHOD = "S256"
	OIDC_PROMPT                = "none"
	OIDC_CODE_CHALLENGE_LENGTH = 43
	OIDC_TOKEN_URL             = "https://www.vodafone.de/mint/oidc/token"
	OIDC_GRANT_TYPE            = "authorization_code"
	// Hard coded API key in the website of vodafone. This is not a secret. It is used to authenticate the API.
	// It is highly likely that this can change and needs to be updated after some time.
	DEFAULT_API_KEY = "aEIoMCae0A933wBL0bLlS6SwSBfkKwM5"
)

var (
	ErrLoginFailed = errors.New("Login failed")
	loginClient    = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

type LoginData struct {
	Cookies      []*http.Cookie
	OIDCResponse *OIDCResponse
}

func (loginData *LoginData) GetToken() string {
	if loginData.OIDCResponse == nil {
		return ""
	}
	return fmt.Sprintf("%s %s", loginData.OIDCResponse.TokenType, loginData.OIDCResponse.AccessToken)
}

func (loginData *LoginData) AuthenticateAPI(request *http.Request) *http.Request {
	loginToken := loginData.GetToken()
	request.Header.Set("Authorization", loginToken)
	request.Header.Set("x-api-key", DEFAULT_API_KEY)
	return request
}

func setDefaultHeaders(request *http.Request) *http.Request {
	request.Header.Set("Accept", DEFAULT_ACCEPT)
	request.Header.Set("Accept-Language", DEFAULT_ACCEPT_LANGUAGE)
	request.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	return request
}

func FromEnvironmentVariables() (*LoginData, error) {
	username := os.Getenv("VODAFONE_USERNAME")
	password := os.Getenv("VODAFONE_PASSWORD")
	credentials, err := Do(username, password)
	return credentials, err
}

func Do(user string, credential string) (*LoginData, error) {
	cookies, err := getMintCookies()
	if err != nil {
		return nil, err
	}
	postBody, err := json.Marshal(map[string]interface{}{
		"authnIdentifier": user,
		"context":         "",
		"conversation":    "",
		"credential":      credential,
		"targetURL":       DEFAULT_TARGET_URL,
	})
	if err != nil {
		return nil, err
	}
	requestBody := bytes.NewBuffer(postBody)
	request, _ := http.NewRequest("POST", DEFAULT_LOGIN_URL, requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", DEFAULT_REFERER)

	request = setDefaultHeaders(request)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(response.Body)
		fmt.Println(string(message))
		return nil, ErrLoginFailed
	}

	cookiesByName := map[string]*http.Cookie{}
	for _, cookie := range cookies {
		cookiesByName[cookie.Name] = cookie
	}
	for _, cookie := range response.Cookies() {
		cookiesByName[cookie.Name] = cookie
	}
	cookies = []*http.Cookie{}
	for _, cookie := range cookiesByName {
		cookies = append(cookies, cookie)
	}

	return oidcLogin(cookies)
}

func getMintCookies() ([]*http.Cookie, error) {

	challenge := randCodeChallenge()
	authorizeURL := oidcAuthorizationURL(challenge)

	request, _ := http.NewRequest("GET", authorizeURL.String(), nil)
	request.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	response, err := loginClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusFound {
		return nil, ErrLoginFailed
	}

	return response.Cookies(), nil
}

func randCodeChallenge() string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, OIDC_CODE_CHALLENGE_LENGTH)
	for i := range b {
		intRand, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		b[i] = letters[intRand.Int64()]
	}
	return string(b)
}

func oidcAuthorizationURL(challenge string) *url.URL {
	authorizeURL, err := url.Parse(OIDC_AUTHORIZE_URL)
	if err != nil {
		panic(err)
	}
	query := authorizeURL.Query()
	hash := sha256.New()
	hash.Write([]byte(challenge))
	hashString := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

	query.Add("response_type", OIDC_RESPONSE_TYPE)
	query.Add("client_id", OIDC_CLIENT_ID)
	query.Add("scope", OIDC_SCOPES)
	query.Add("redirect_uri", OIDC_REDIRECT_URI)
	query.Add("code_challenge", hashString)
	query.Add("code_challenge_method", OIDC_CODE_CHALLENGE_METHOD)
	query.Add("prompt", OIDC_PROMPT)
	authorizeURL.RawQuery = query.Encode()

	return authorizeURL
}

func oidcLogin(cookies []*http.Cookie) (*LoginData, error) {
	challenge := randCodeChallenge()
	authorizeURL := oidcAuthorizationURL(challenge)
	query := authorizeURL.Query()
	query.Add("state", "")
	authorizeURL.RawQuery = query.Encode()

	request, _ := http.NewRequest("GET", authorizeURL.String(), nil)
	setDefaultHeaders(request)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	response, err := loginClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusFound {
		return nil, ErrLoginFailed
	}
	location := response.Header.Get("Location")
	if len(location) == 0 {
		return nil, ErrLoginFailed
	}

	urlLocation, err := url.Parse(location)
	if err != nil {
		return nil, ErrLoginFailed
	}

	code := urlLocation.Query().Get("code")
	if len(code) == 0 {
		return nil, ErrLoginFailed
	}

	oidcResponse, err := oidcToken(code, challenge, cookies)
	if err != nil {
		return nil, err
	}

	loginData := &LoginData{
		Cookies:      cookies,
		OIDCResponse: oidcResponse,
	}
	return loginData, nil
}

type OIDCResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	IDToken     string `json:"id_token"`
}

func oidcToken(code string, challenge string, cookies []*http.Cookie) (*OIDCResponse, error) {
	tokenURL, err := url.Parse(OIDC_TOKEN_URL)
	if err != nil {
		panic(err)
	}
	query := tokenURL.Query()
	query.Set("client_id", OIDC_CLIENT_ID)
	query.Set("grant_type", OIDC_GRANT_TYPE)
	query.Set("code", code)
	query.Set("code_verifier", challenge)
	query.Set("redirect_uri", OIDC_REDIRECT_URI)
	tokenURL.RawQuery = query.Encode()

	request, err := http.NewRequest("POST", tokenURL.String(), nil)
	if err != nil {
		panic(err)
	}
	setDefaultHeaders(request)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, ErrLoginFailed
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseString, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(responseString))

		return nil, ErrLoginFailed
	}
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, ErrLoginFailed
	}

	oidcResponse := &OIDCResponse{}
	err = json.Unmarshal(payload, oidcResponse)
	if err != nil {
		return nil, err
	}

	return oidcResponse, nil
}
