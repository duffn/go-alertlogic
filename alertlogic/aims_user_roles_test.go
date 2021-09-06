package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	getAssignedRolesPath   = fmt.Sprintf("/%s/%s/users/%s/roles", aimsServicePath, testAccountId, testUserId)
	getAssignedRoleIDsPath = fmt.Sprintf("/%s/%s/users/%s/role_ids", aimsServicePath, testAccountId, testUserId)
	getUserPermissionsPath = fmt.Sprintf("/%s/%s/users/%s/permissions", aimsServicePath, testAccountId, testUserId)
)

func TestAims_GetAssignedRoles(t *testing.T) {
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

	mux.HandleFunc(getAssignedRolesPath, func(w http.ResponseWriter, r *http.Request) {
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

	roles, err := client.GetAssignedRoles(testUserId)

	if assert.NoError(t, err) {
		assert.Equal(t, roles, want)
	}
}

func TestAims_GetAssignedRolesMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getAssignedRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetAssignedRoles(testUserId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_GetAssignedRolesUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getAssignedRolesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetAssignedRoles(testUserId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_GetAssignedRoleIDs(t *testing.T) {
	setup()
	defer teardown()

	const response = `{"role_ids": ["F578CCE5-9574-4489-BF05-A04075838DE3", "2A33175D-86EF-44B5-AA39-C9549F6306DF"]}`

	mux.HandleFunc(getAssignedRoleIDsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := RoleIdsList{RoleIds: []string{"F578CCE5-9574-4489-BF05-A04075838DE3", "2A33175D-86EF-44B5-AA39-C9549F6306DF"}}

	roleIds, err := client.GetAssignedRoleIDs(testUserId)

	if assert.NoError(t, err) {
		assert.Equal(t, roleIds, want)
	}
}

func TestAims_GetAssignedRoleIDsMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getAssignedRoleIDsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetAssignedRoleIDs(testUserId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_GetAssignedRoleIDsUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getAssignedRoleIDsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetAssignedRoleIDs(testUserId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestAims_GetUserPermissions(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"permissions": [
			{"*:managed:*:*": "allowed"},
			{"aims:own:update:role": "denied"},
			{"aims:own:delete:role": "denied"}
		]
	}`

	mux.HandleFunc(getUserPermissionsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := PermissionsList{
		Permissions: []map[string]Permission{
			{"*:managed:*:*": Allowed},
			{"aims:own:update:role": Denied},
			{"aims:own:delete:role": Denied},
		},
	}

	perms, err := client.GetUserPermissions(testUserId)

	if assert.NoError(t, err) {
		assert.Equal(t, perms, want)
	}
}

func TestAims_GetUserPermissionsMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getUserPermissionsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetUserPermissions(testUserId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAims_GetUserPermissionsUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getUserPermissionsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetUserPermissions(testUserId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}
