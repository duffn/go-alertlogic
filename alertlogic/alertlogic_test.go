package alertlogic

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// mux is the HTTP request multiplexer used with the test server
	mux *http.ServeMux
	// client is the API client being tested
	client *API
	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
	// testUnmarshalError is a JSON unmarshaling error used for testing
	testUnmarshalError = fmt.Sprintf("%s: invalid character 'o' in literal null (expecting 'u')", errUnmarshalError)
	// testNotFoundError is a not found error used for testing
	testNotFoundError = fmt.Sprintf("%s: HTTP status 404: content \"\"", errMakeRequestError)
	// testAccountId is an Alert Logic account ID used for testing
	testAccountId = "12345678"
)

const (
	testRelatedAccountId = "98765432"
	testUserId           = "715A4EC0-9833-4D6E-9C03-A537E3F98D23"
	testEmail            = "bob@bobloblawlaw.com"
	testUserFullName     = "Bob Loblaw"
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// API client configured to use test server
	client, _ = NewWithApiToken(testAccountId, "my_token")
	client.BaseURL = server.URL
}

func teardown() {
	server.Close()
}

var credentials = []struct {
	username string
	password string
}{
	{"", ""},
	{"username", ""},
	{"", "password"},
}

func TestAlertLogic_NewMissingCredentials(t *testing.T) {
	for _, tt := range credentials {
		_, err := NewWithUsernameAndPassword(testAccountId, tt.username, tt.password)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errEmptyUsernameOrPassword)
	}
}

func TestAlertLogic_NewMissingApiToken(t *testing.T) {
	_, err := NewWithApiToken(testAccountId, "")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errEmptyApiToken)
}

func TestAlertLogic_NewMissingAccountId(t *testing.T) {
	_, err := NewWithUsernameAndPassword("", "username", "password")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errEmptyAccountId)

	_, errApiToken := NewWithApiToken("", "abcd1234")
	assert.Error(t, errApiToken)
	assert.Equal(t, errApiToken.Error(), errEmptyAccountId)
}
