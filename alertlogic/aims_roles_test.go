package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	listRolesPath            = fmt.Sprintf("/%s/%s/roles", aimsServicePath, testAccountId)
	listGlobalRolesPath      = fmt.Sprintf("/%s/roles", aimsServicePath)
	getRoleDetailsPath       = fmt.Sprintf("/%s/%s/roles/%s", aimsServicePath, testAccountId, testRoleId)
	getGlobalRoleDetailsPath = fmt.Sprintf("/%s/roles/%s", aimsServicePath, testRoleId)
)

func TestAims_ListRoles(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"roles": [{
			"id": "F578CCE5-9574-4489-BF05-A04075838DE3",
			"account_id": "12345678",
			"name": "Read Only",
			"permissions": {
				"*:own:list:*": "allowed",
				"*:own:get:*": "allowed"
			},
			"legacy_permissions": [
				"PERM1",
				"PERM2"
			],
			"version": 1,
			"global": false,
			"created": {
				"at": 1430184599,
				"by": "System"
			},
			"modified": {
				"at": 1430184599,
				"by": "System"
			}
		},
		{
			"id": "2A33175D-86EF-44B5-AA39-C9549F6306DF",
			"account_id": "12345678",
			"name": "Power User",
			"permissions": {
				"aims:own:create:*": "denied",
				"*:own:*:*": "allowed"
			},
			"legacy_permissions": [],
			"version": 1,
			"global": false,
			"created": {
				"at": 1430184599,
				"by": "System"
			},
			"modified": {
				"at": 1430184599,
				"by": "System"
			}
		}]
	}`

	mux.HandleFunc(listRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := RolesList{
		Roles: []Role{
			{
				ID:                testRoleId,
				AccountID:         testAccountId,
				Name:              "Read Only",
				Permissions:       map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed},
				LegacyPermissions: []string{"PERM1", "PERM2"},
				Version:           1,
				Global:            false,
				Created:           ModifiedCreated{At: 1430184599, By: "System"},
				Modified:          ModifiedCreated{At: 1430184599, By: "System"},
			},
			{
				ID:                "2A33175D-86EF-44B5-AA39-C9549F6306DF",
				AccountID:         testAccountId,
				Name:              "Power User",
				Permissions:       map[string]Permission{"aims:own:create:*": Denied, "*:own:*:*": Allowed},
				Version:           1,
				Global:            false,
				LegacyPermissions: []string{},
				Created:           ModifiedCreated{At: 1430184599, By: "System"},
				Modified:          ModifiedCreated{At: 1430184599, By: "System"},
			},
		},
	}

	roles, err := client.ListRoles()

	if assert.NoError(t, err) {
		assert.Equal(t, roles, want)
	}
}

func TestAims_ListRolesMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.ListRoles()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_ListRolesUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.ListRoles()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_ListGlobalRoles(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"roles": [{
			"id": "F578CCE5-9574-4489-BF05-A04075838DE3",
			"account_id": "*",
			"name": "Read Only",
			"permissions": {
				"*:own:list:*": "allowed",
				"*:own:get:*": "allowed"
			},
			"legacy_permissions": [
				"PERM1",
				"PERM2"
			],
			"version": 1,
			"global": true,
			"created": {
				"at": 1430184599,
				"by": "System"
			},
			"modified": {
				"at": 1430184599,
				"by": "System"
			}
		},
		{
			"id": "2A33175D-86EF-44B5-AA39-C9549F6306DF",
			"account_id": "*",
			"name": "Power User",
			"permissions": {
				"aims:own:create:*": "denied",
				"*:own:*:*": "allowed"
			},
			"legacy_permissions": [],
			"version": 1,
			"global": true,
			"created": {
				"at": 1430184599,
				"by": "System"
			},
			"modified": {
				"at": 1430184599,
				"by": "System"
			}
		}]
	}`

	mux.HandleFunc(listGlobalRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := RolesList{
		Roles: []Role{
			{
				ID:                testRoleId,
				AccountID:         "*",
				Name:              "Read Only",
				Permissions:       map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed},
				LegacyPermissions: []string{"PERM1", "PERM2"},
				Version:           1,
				Global:            true,
				Created:           ModifiedCreated{At: 1430184599, By: "System"},
				Modified:          ModifiedCreated{At: 1430184599, By: "System"},
			},
			{
				ID:                "2A33175D-86EF-44B5-AA39-C9549F6306DF",
				AccountID:         "*",
				Name:              "Power User",
				Permissions:       map[string]Permission{"aims:own:create:*": Denied, "*:own:*:*": Allowed},
				Version:           1,
				Global:            true,
				LegacyPermissions: []string{},
				Created:           ModifiedCreated{At: 1430184599, By: "System"},
				Modified:          ModifiedCreated{At: 1430184599, By: "System"},
			},
		},
	}

	roles, err := client.ListGlobalRoles()

	if assert.NoError(t, err) {
		assert.Equal(t, roles, want)
	}
}

func TestAims_ListGlobalRolesMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listGlobalRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.ListGlobalRoles()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_ListGlobalRolesUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listGlobalRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.ListGlobalRoles()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_GetRoleDetails(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"id": "F578CCE5-9574-4489-BF05-A04075838DE3",
		"account_id": "12345678",
		"name": "Read Only",
		"permissions": {
			"*:own:list:*": "allowed",
			"*:own:get:*": "allowed"
		},
		"legacy_permissions": [
			"PERM1",
			"PERM2"
		],
		"version": 1,
		"global": false,
		"created": {
			"at": 1430184599,
			"by": "System"
		},
		"modified": {
			"at": 1430184599,
			"by": "System"
		}
		
	}`

	mux.HandleFunc(getRoleDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := Role{
		ID:                testRoleId,
		AccountID:         testAccountId,
		Name:              "Read Only",
		Permissions:       map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed},
		LegacyPermissions: []string{"PERM1", "PERM2"},
		Version:           1,
		Global:            false,
		Created:           ModifiedCreated{At: 1430184599, By: "System"},
		Modified:          ModifiedCreated{At: 1430184599, By: "System"},
	}

	roles, err := client.GetRoleDetails(testRoleId)

	if assert.NoError(t, err) {
		assert.Equal(t, roles, want)
	}
}

func TestAims_GetRoleDetailsMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getRoleDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetRoleDetails(testRoleId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_GetRoleDetailsUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getRoleDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetRoleDetails(testRoleId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_GetGlobalRoleDetails(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"id": "F578CCE5-9574-4489-BF05-A04075838DE3",
		"account_id": "*",
		"name": "Read Only",
		"permissions": {
			"*:own:list:*": "allowed",
			"*:own:get:*": "allowed"
		},
		"legacy_permissions": [
			"PERM1",
			"PERM2"
		],
		"version": 1,
		"global": true,
		"created": {
			"at": 1430184599,
			"by": "System"
		},
		"modified": {
			"at": 1430184599,
			"by": "System"
		}
		
	}`

	mux.HandleFunc(getGlobalRoleDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := Role{
		ID:                testRoleId,
		AccountID:         "*",
		Name:              "Read Only",
		Permissions:       map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed},
		LegacyPermissions: []string{"PERM1", "PERM2"},
		Version:           1,
		Global:            true,
		Created:           ModifiedCreated{At: 1430184599, By: "System"},
		Modified:          ModifiedCreated{At: 1430184599, By: "System"},
	}

	roles, err := client.GetGlobalRoleDetails(testRoleId)

	if assert.NoError(t, err) {
		assert.Equal(t, roles, want)
	}
}

func TestAims_GetGlobalRoleDetailsMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getGlobalRoleDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetGlobalRoleDetails(testRoleId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_GetGlobalRoleDetailsUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getGlobalRoleDetailsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetGlobalRoleDetails(testRoleId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}
