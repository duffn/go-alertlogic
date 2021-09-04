package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	authenticatePath       = fmt.Sprintf("/%s/authenticate", aimsServicePath)
	createUserPath         = fmt.Sprintf("/%s/%s/users", aimsServicePath, testAccountId)
	deleteUserPath         = fmt.Sprintf("/%s/%s/users/%s", aimsServicePath, testAccountId, testUserId)
	listUsersByEmailPath   = fmt.Sprintf("/%s/users/email/%s", aimsServicePath, testEmail)
	getUserDetailsByIdPath = fmt.Sprintf("/%s/user/%s", aimsServicePath, testUserId)
)

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
				"email": "bob@bobloblawlaw.com",
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
				Name:      testUserFullName,
				Email:     testEmail,
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

	var mobilePhone string = "123-555-0123"

	want := User{
		ID:          testUserId,
		AccountID:   testAccountId,
		Name:        testUserFullName,
		Email:       testEmail,
		Username:    testEmail,
		Active:      true,
		Version:     1,
		MobilePhone: &mobilePhone,
		Locked:      false,
		LinkedUsers: []LinkedUser{},
		Created:     ModifiedCreated{At: 1430185015, By: "System"},
		Modified:    ModifiedCreated{At: 1430185015, By: "System"},
	}

	user, err := client.CreateUser(CreateUserRequest{Email: testEmail, Name: testUserFullName}, false)

	if assert.NoError(t, err) {
		assert.Equal(t, user, want)
	}

	user, err = client.CreateUser(CreateUserRequest{Email: testEmail, Name: testUserFullName, Password: "password"}, true)

	if assert.NoError(t, err) {
		assert.Equal(t, user, want)
	}
}

func TestAims_CreateUserOneTimePasswordMissingPassword(t *testing.T) {
	_, err := client.CreateUser(CreateUserRequest{Email: testEmail, Name: testUserFullName}, true)

	assert.Error(t, err, "oneTimePassword must be accompanied by CreateUserRequest.Password")
}

func TestAims_CreateUserMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(createUserPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.CreateUser(CreateUserRequest{Email: testEmail, Name: "Bob Loblaw"}, false)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_CreateUserUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(createUserPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.CreateUser(CreateUserRequest{Email: testEmail, Name: testUserFullName}, false)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_DeleteUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(deleteUserPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	deleteUserResponse, err := client.DeleteUser(testUserId)

	if assert.NoError(t, err) {
		assert.Equal(t, deleteUserResponse, http.StatusNoContent)
	}
}

func TestAims_DeleteUserError(t *testing.T) {
	setup()
	defer teardown()

	errorResponse := `{"error":"self_delete_error"}`

	mux.HandleFunc(deleteUserPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, errorResponse)
	})

	respCode, err := client.DeleteUser(testUserId)

	assert.Error(t, err)
	assert.Equal(t, respCode, http.StatusBadRequest)
	assert.Equal(t, err.Error(), fmt.Sprintf("error from makeRequest: %s", errorResponse))
}

func TestAims_ListUsersByEmail(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"users": [{
			"id": "715A4EC0-9833-4D6E-9C03-A537E3F98D23",
			"account_id": "12345678",
			"name": "Bob Loblaw",
			"username": "bob@bobloblawlaw.com",
			"email": "bob@bobloblawlaw.com",
			"active": true,
			"locked": false,
			"version": 1,
			"linked_users": [],
			"mfa_enabled": true,
			"created": {
				"at": 1430185015,
				"by": "System"
			},
			"modified": {
				"at": 1430185015,
				"by": "System"
			}
		}]
	}`

	mux.HandleFunc(listUsersByEmailPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	var mfaEnabled bool = true

	want := UserList{
		Users: []User{
			{
				ID:          testUserId,
				AccountID:   testAccountId,
				Name:        testUserFullName,
				Email:       testEmail,
				Username:    testEmail,
				Active:      true,
				Version:     1,
				MfaEnabled:  &mfaEnabled,
				Locked:      false,
				LinkedUsers: []LinkedUser{},
				Created:     ModifiedCreated{At: 1430185015, By: "System"},
				Modified:    ModifiedCreated{At: 1430185015, By: "System"},
			},
		},
	}

	users, err := client.ListUsersByEmail(testEmail)

	if assert.NoError(t, err) {
		assert.Equal(t, users, want)
	}
}

func TestAims_ListUsersByEmailMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listUsersByEmailPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.ListUsersByEmail(testEmail)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_ListUsersByEmailUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listUsersByEmailPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.ListUsersByEmail(testEmail)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_GetUserDetailsById(t *testing.T) {
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
		"role_ids": ["2A33175D-86EF-44B5-AA39-C9549F6306DF"],
		"user_credential": {
			"version": 2,
			"one_time_password": false,
			"last_login": 1548880711,
			"created": {
				"at": 1430185015,
				"by": "System"
			},
			"modified": {
				"at": 1430185015,
				"by": "System"
			}
		},
		"access_keys": [{
			"label": "api access",
			"last_login": 0,
			"created": {
				"at": 1525410880,
				"by": "System"
			},
			"modified": {
				"at": 1525410880,
				"by": "System"
			},
			"access_key_id": "61fb235617960503"
		}],
		"created": {
			"at": 1430185015,
			"by": "System"
		},
		"modified": {
			"at": 1430185015,
			"by": "System"
		}
	}`

	mux.HandleFunc(getUserDetailsByIdPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	var mobilePhone string = "123-555-0123"
	var userCredential UserCredential = UserCredential{
		Version:         2,
		OneTimePassword: false,
		LastLogin:       1548880711,
		Created:         ModifiedCreated{At: 1430185015, By: "System"},
		Modified:        ModifiedCreated{At: 1430185015, By: "System"},
	}
	var accessKeys []AccessKey = []AccessKey{
		{
			Label:       "api access",
			LastLogin:   0,
			Created:     ModifiedCreated{At: 1525410880, By: "System"},
			Modified:    ModifiedCreated{At: 1525410880, By: "System"},
			AccessKeyID: "61fb235617960503",
		},
	}

	want := User{
		ID:             testUserId,
		AccountID:      testAccountId,
		Name:           testUserFullName,
		Email:          testEmail,
		Username:       testEmail,
		Active:         true,
		Version:        1,
		MfaEnabled:     nil,
		MobilePhone:    &mobilePhone,
		Locked:         false,
		LinkedUsers:    []LinkedUser{},
		Created:        ModifiedCreated{At: 1430185015, By: "System"},
		Modified:       ModifiedCreated{At: 1430185015, By: "System"},
		UserCredential: &userCredential,
		AccessKeys:     &accessKeys,
		RoleIds:        &[]string{"2A33175D-86EF-44B5-AA39-C9549F6306DF"},
	}

	user, err := client.GetUserDetailsById(testUserId, true, true, true)

	if assert.NoError(t, err) {
		assert.Equal(t, user, want)
	}
}

func TestAims_GetUserDetailsByIdWithNoParams(t *testing.T) {
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

	mux.HandleFunc(getUserDetailsByIdPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	var mobilePhone string = "123-555-0123"

	want := User{
		ID:          testUserId,
		AccountID:   testAccountId,
		Name:        testUserFullName,
		Email:       testEmail,
		Username:    testEmail,
		Active:      true,
		Version:     1,
		MfaEnabled:  nil,
		MobilePhone: &mobilePhone,
		Locked:      false,
		LinkedUsers: []LinkedUser{},
		Created:     ModifiedCreated{At: 1430185015, By: "System"},
		Modified:    ModifiedCreated{At: 1430185015, By: "System"},
	}

	user, err := client.GetUserDetailsById(testUserId, false, false, false)

	if assert.NoError(t, err) {
		assert.Equal(t, user, want)
	}
}

func TestAims_GetUserDetailsByIdMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getUserDetailsByIdPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetUserDetailsById(testUserId, true, true, true)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_GetUserDetailsByIdUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getUserDetailsByIdPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetUserDetailsById(testUserId, true, true, true)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}
