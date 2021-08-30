package alertlogic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// apiURL is the base URL for the API.
const apiURL = "https://api.cloudinsight.alertlogic.com"

// API holds the configuration for the current API client.
type API struct {
	Username   string
	Password   string
	APIToken   string
	BaseURL    string
	AccountID  string
	UserAgent  string
	headers    http.Header
	httpClient *http.Client
}

// ModifiedCreated holds the created or modified response from the API.
type ModifiedCreated struct {
	At int    `json:"at,omitempty"`
	By string `json:"by,omitempty"`
}

// newClient creates a new API client.
func newClient(accountId string) (*API, error) {
	if accountId == "" {
		return nil, errors.New(errEmptyAccountId)
	}

	api := &API{
		BaseURL:   apiURL,
		AccountID: accountId,
		headers:   make(http.Header),
	}

	if api.httpClient == nil {
		api.httpClient = http.DefaultClient
	}

	return api, nil
}

// NewWithApiToken creates a new Alert Logic API client using an API token.
func NewWithApiToken(accountId string, apiToken string) (*API, error) {
	if apiToken == "" {
		return nil, errors.New(errEmptyApiToken)
	}

	api, err := newClient(accountId)
	if err != nil {
		return nil, err
	}

	api.APIToken = apiToken
	return api, nil
}

// NewWithUsernameAndPassword creates a new Alert Logic API client using a username and password.
// Username and password should be an access key and secret key as noted in the documentation
// https://docs.alertlogic.com/prepare/access-key-management.htm, but can also be your Alert Logic
// UI email and password.
func NewWithUsernameAndPassword(accountId string, username string, password string) (*API, error) {
	if username == "" || password == "" {
		return nil, errors.New(errEmptyUsernameOrPassword)
	}

	api, err := newClient(accountId)
	if err != nil {
		return nil, err
	}

	api.Username = username
	api.Password = password

	authenticateResponse, err := api.Authenticate()
	if err != nil {
		return nil, err
	}

	api.APIToken = authenticateResponse.Authentication.Token
	return api, nil
}

// makeRequest makes an HTTP request.
func (api *API) makeRequest(
	method string,
	path string,
	headers http.Header,
	params map[string]string,
	body interface{},
) ([]byte, int, error) {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, 0, errors.Wrap(err, "error marshalling body to JSON")
		}
	}

	var requestBody io.Reader
	if jsonBody != nil {
		requestBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", api.BaseURL, path),
		requestBody,
	)
	if err != nil {
		return nil, 0, errors.Wrap(err, errMakeRequestError)
	}

	combinedHeaders := make(http.Header)
	copyHeader(combinedHeaders, api.headers)
	copyHeader(combinedHeaders, headers)
	req.Header = combinedHeaders

	if api.Username != "" && api.Password != "" {
		req.SetBasicAuth(api.Username, api.Password)
	}

	if api.APIToken != "" {
		req.Header.Set("X-Aims-Auth-Token", api.APIToken)
	}

	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := api.httpClient.Do(req)

	if err != nil {
		return nil, 0, errors.Wrap(err, errMakeRequestError)
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	switch {
	case resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices:
	case resp.StatusCode == http.StatusUnauthorized:
		return nil, resp.StatusCode, errors.Errorf("HTTP status %d: invalid credentials", resp.StatusCode)
	case resp.StatusCode == http.StatusForbidden:
		return nil, resp.StatusCode, errors.Errorf("HTTP status %d: insufficient permissions", resp.StatusCode)
	case resp.StatusCode == http.StatusServiceUnavailable,
		resp.StatusCode == http.StatusBadGateway,
		resp.StatusCode == http.StatusGatewayTimeout,
		resp.StatusCode == 522,
		resp.StatusCode == 523,
		resp.StatusCode == 524:
		return nil, resp.StatusCode, errors.Errorf("HTTP status %d: service failure", resp.StatusCode)
	case resp.StatusCode == 400:
		return nil, resp.StatusCode, errors.Errorf("%s", respBody)
	default:
		var s string
		if respBody != nil {
			s = string(respBody)
		}
		return nil, resp.StatusCode, errors.Errorf("HTTP status %d: content %q", resp.StatusCode, s)
	}

	return respBody, resp.StatusCode, nil
}

// copyHeader copies all headers for `source` and sets them on `target`.
func copyHeader(target http.Header, source http.Header) {
	for k, vs := range source {
		target[k] = vs
	}
}
