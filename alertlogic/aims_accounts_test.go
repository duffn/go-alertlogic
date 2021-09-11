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

func TestAims_UpdateAccountDetails(t *testing.T) {
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
		"mfa_required": false,
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
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := AccountDetails{
		ID:                  "12345678",
		Name:                "Company Name",
		Active:              true,
		MfaRequired:         false,
		Version:             1,
		AccessibleLocations: []string{"insight-us-virginia"},
		DefaultLocation:     "insight-us-virginia",
		Created:             ModifiedCreated{At: 1430184599, By: "System"},
		Modified:            ModifiedCreated{At: 1430184599, By: "System"},
	}

	user, err := client.UpdateAccountDetails(UpdateAccountDetailsRequest{MfaRequired: false})

	if assert.NoError(t, err) {
		assert.Equal(t, user, want)
	}
}
