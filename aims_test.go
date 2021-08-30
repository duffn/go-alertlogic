package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testRelatedAccountId = "98765432"

var accountDetailsPath = fmt.Sprintf("/%s/%s/account", aimsServicePath, testAccountId)
var authenticatePath = fmt.Sprintf("/%s/authenticate", aimsServicePath)
var accountRelationshipPath = fmt.Sprintf("/%s/%s/accounts/%s/%s", aimsServicePath, testAccountId, Managed, testRelatedAccountId)
var createUserPath = fmt.Sprintf("/%s/%s/users", aimsServicePath, testAccountId)

func TestAims_Authenticate(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"authentication": {
			"user": {
				"id": "715A4EC0-9833-4D6E-9C03-A537E3F98D23",
				"account_id": "12345678",
				"name": "Bob Loblaw",
				"email": "bob@loblawlaw.com",
				"active": true,
				"locked": false,
				"version": 1,
				"created": {
					"at": 1430183768,
					"by": "System"
				},
				"modified": {
					"at": 1430183768,
					"by": "System"
				}
			},
			"account": {
				"id": "12345678",
				"name": "Loblaw Law",
				"active": true,
				"version": 1,
				"accessible_locations": ["insight-us-virginia"],
				"default_location": "insight-us-virginia",
				"created": {
					"by": "System",
					"at": 1436482061
				},
				"modified": {
					"by": "System",
					"at": 1436482061
				}
			},
			"token": "my_long_token",
			"token_expiration": 1434042731
		}
	}`

	mux.HandleFunc(authenticatePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := AuthenticateResponse{
		Authentication{
			Token:           "my_long_token",
			TokenExpiration: 1434042731,
			User: User{
				ID:        "715A4EC0-9833-4D6E-9C03-A537E3F98D23",
				AccountID: testAccountId,
				Name:      "Bob Loblaw",
				Email:     "bob@loblawlaw.com",
				Active:    true,
				Locked:    false,
				Version:   1,
				Created: ModifiedCreated{
					At: 1430183768,
					By: "System",
				},
				Modified: ModifiedCreated{
					At: 1430183768,
					By: "System",
				},
			},
			Account: Account{
				ID:                  testAccountId,
				Name:                "Loblaw Law",
				Active:              true,
				Version:             1,
				AccessibleLocations: []string{"insight-us-virginia"},
				DefaultLocation:     "insight-us-virginia",
				Created:             ModifiedCreated{At: 1436482061, By: "System"},
				Modified:            ModifiedCreated{At: 1436482061, By: "System"},
			},
		},
	}

	authenticateResponse, err := client.Authenticate()
	if assert.NoError(t, err) {
		assert.Equal(t, authenticateResponse, want)
	}
}

func TestAims_AuthenticateMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(authenticatePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.Authenticate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_AuthenticateUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(authenticatePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.Authenticate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

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

	want := AccountDetailsResponse{
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

func TestAims_CreateUser(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"id": "715A4EC0-9833-4D6E-9C03-A537E3F98D23",
		"account_id": "12345678",
		"name": "Bob Loblaw",
		"username": "bob@bobloblawlaw.com",
		"email": "bob@bobloblawlaw.com",
		"active": true,
		"locked": false,
		"version": 1,
		"linked_users": [],
		"mobile_phone": "123-555-0123",
		"created": {
			"at": 1430185015,
			"by": "System"
		},
		"modified": {
			"at": 1430185015,
			"by": "System"
		}
	}`

	mux.HandleFunc(createUserPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := CreateUserResponse{
		ID:          "715A4EC0-9833-4D6E-9C03-A537E3F98D23",
		AccountID:   testAccountId,
		Name:        "Bob Loblaw",
		Email:       "bob@bobloblawlaw.com",
		Username:    "bob@bobloblawlaw.com",
		Active:      true,
		Version:     1,
		MobilePhone: "123-555-0123",
		Locked:      false,
		LinkedUsers: []LinkedUser{},
		Created:     ModifiedCreated{At: 1430185015, By: "System"},
		Modified:    ModifiedCreated{At: 1430185015, By: "System"},
	}

	user, err := client.CreateUser(CreateUserRequest{Email: "bob@bobloblawlaw.com", Name: "Bob Loblaw"}, false)

	if assert.NoError(t, err) {
		assert.Equal(t, user, want)
	}

	user, err = client.CreateUser(CreateUserRequest{Email: "bob@bobloblawlaw.com", Name: "Bob Loblaw", Password: "password"}, true)

	if assert.NoError(t, err) {
		assert.Equal(t, user, want)
	}
}

func TestAims_CreateUserOneTimePasswordMissingPassword(t *testing.T) {
	_, err := client.CreateUser(CreateUserRequest{Email: "bob@bobloblawlaw.com", Name: "Bob Loblaw"}, true)

	assert.Error(t, err, "oneTimePassword must be accompanied by CreateUserRequest.Password")
}
