package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	accountDetailsPath      = fmt.Sprintf("/%s/%s/account", aimsServicePath, testAccountId)
	accountRelationshipPath = fmt.Sprintf("/%s/%s/accounts/%s/%s", aimsServicePath, testAccountId, Managed, testRelatedAccountId)
)

func TestAims_GetAccountDetails(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"id": "12345678",
		"name": "Company Name",
		"active": true,
		"version": 1,
		"accessible_locations": ["insight-us-virginia"],
		"default_location": "insight-us-virginia",
		"created": {
			"at": 1430184599,
			"by": "System"
		},
		"modified": {
			"at": 1430184599,
			"by": "System"
		}
	}`

	mux.HandleFunc(accountDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := AccountDetails{
		ID:                  "12345678",
		Name:                "Company Name",
		Active:              true,
		Version:             1,
		AccessibleLocations: []string{"insight-us-virginia"},
		DefaultLocation:     "insight-us-virginia",
		Created:             ModifiedCreated{At: 1430184599, By: "System"},
		Modified:            ModifiedCreated{At: 1430184599, By: "System"},
	}

	accountDetails, err := client.GetAccountDetails()

	if assert.NoError(t, err) {
		assert.Equal(t, accountDetails, want)
	}
}

func TestAims_GetAccountDetailsMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(accountDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetAccountDetails()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_GetAccountDetailsUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(accountDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetAccountDetails()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_GetAccountRelationship(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(accountRelationshipPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	accountDetailsResponse, err := client.GetAccountRelationship(testRelatedAccountId, Managed)

	if assert.NoError(t, err) {
		assert.Equal(t, accountDetailsResponse, http.StatusNoContent)
	}
}

func TestAims_GetAccountRelationshipNotFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(accountRelationshipPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.GetAccountRelationship(testRelatedAccountId, Managed)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testNotFoundError)
}
