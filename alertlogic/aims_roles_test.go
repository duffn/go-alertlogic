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
	deleteRolePath           = fmt.Sprintf("/%s/%s/roles/%s", aimsServicePath, testAccountId, testRoleId)
	updateRolePath           = fmt.Sprintf("/%s/%s/roles/%s", aimsServicePath, testAccountId, testRoleId)
	createRolePath           = fmt.Sprintf("/%s/%s/roles", aimsServicePath, testAccountId)
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

func TestAims_UpdateRoleDetails(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"id": "F578CCE5-9574-4489-BF05-A04075838DE3",
		"account_id": "12345678",
		"name": "Read Only",
		"permissions": {
			"*:own:list:*": "allowed",
			"*:own:get:*": "allowed",
			"*:own:*:*": "allowed"
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

	mux.HandleFunc(updateRolePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := Role{
		ID:                testRoleId,
		AccountID:         testAccountId,
		Name:              "Read Only",
		Permissions:       map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed, "*:own:*:*": Allowed},
		LegacyPermissions: []string{"PERM1", "PERM2"},
		Version:           1,
		Global:            false,
		Created:           ModifiedCreated{At: 1430184599, By: "System"},
		Modified:          ModifiedCreated{At: 1430184599, By: "System"},
	}

	role, err := client.UpdateRoleDetails(testRoleId, UpdateRoleRequest{Permissions: map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed, "*:own:*:*": Allowed}})

	if assert.NoError(t, err) {
		assert.Equal(t, role, want)
	}
}

func TestAims_DeleteRole(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(deleteRolePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	deleteRoleResponse, err := client.DeleteRole(testRoleId)

	if assert.NoError(t, err) {
		assert.Equal(t, deleteRoleResponse, http.StatusNoContent)
	}
}

func TestAims_CreateRole(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"id": "F578CCE5-9574-4489-BF05-A04075838DE3",
		"account_id": "12345678",
		"name": "Read Only",
		"permissions": {
			"*:own:list:*": "allowed",
			"*:own:get:*": "allowed",
			"*:own:*:*": "allowed"
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

	mux.HandleFunc(createRolePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := Role{
		ID:                testRoleId,
		AccountID:         testAccountId,
		Name:              "Read Only",
		Permissions:       map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed, "*:own:*:*": Allowed},
		LegacyPermissions: []string{"PERM1", "PERM2"},
		Version:           1,
		Global:            false,
		Created:           ModifiedCreated{At: 1430184599, By: "System"},
		Modified:          ModifiedCreated{At: 1430184599, By: "System"},
	}

	role, err := client.CreateRole(CreateRoleRequest{Name: "Read Only", Permissions: map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed, "*:own:*:*": Allowed}})

	if assert.NoError(t, err) {
		assert.Equal(t, role, want)
	}
}
