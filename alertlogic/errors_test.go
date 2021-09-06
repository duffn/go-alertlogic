package alertlogic

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func callFunction(obj interface{}, fn string, args map[string]interface{}) (res []reflect.Value) {
	method := reflect.ValueOf(obj).MethodByName(fn)
	var inputs []reflect.Value
	for _, v := range args {
		inputs = append(inputs, reflect.ValueOf(v))
	}
	return method.Call(inputs)
}

type RequestError struct {
	Group        string
	Path         string
	Method       string
	FunctionName string
	Arguments    map[string]interface{}
}

var tests = []RequestError{
	{
		Group:        "accounts",
		Path:         accountDetailsPath,
		Method:       "GET",
		FunctionName: "GetAccountDetails",
		Arguments:    nil,
	},
	{
		Group:        "accounts",
		Path:         accountDetailsPath,
		Method:       "POST",
		FunctionName: "UpdateAccountDetails",
		Arguments: map[string]interface{}{
			"updateAccountDetailsRequest": UpdateAccountDetailsRequest{MfaRequired: false},
		},
	},
	{
		Group:        "roles",
		Path:         listRolesPath,
		Method:       "GET",
		FunctionName: "ListRoles",
		Arguments:    nil,
	},
	{
		Group:        "roles",
		Path:         listGlobalRolesPath,
		Method:       "GET",
		FunctionName: "ListGlobalRoles",
		Arguments:    nil,
	},
	{
		Group:        "roles",
		Path:         getRoleDetailsPath,
		Method:       "GET",
		FunctionName: "GetRoleDetails",
		Arguments: map[string]interface{}{
			"roleId": testRoleId,
		},
	},
	{
		Group:        "roles",
		Path:         getGlobalRoleDetailsPath,
		Method:       "GET",
		FunctionName: "GetGlobalRoleDetails",
		Arguments: map[string]interface{}{
			"roleId": testRoleId,
		},
	},
	{
		Group:        "roles",
		Path:         updateRolePath,
		Method:       "POST",
		FunctionName: "UpdateRoleDetails",
		Arguments: map[string]interface{}{
			"roleId": testRoleId,
			"role":   UpdateRoleRequest{Permissions: map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed, "*:own:*:*": Allowed}},
		},
	},
	{
		Group:        "roles",
		Path:         createRolePath,
		Method:       "POST",
		FunctionName: "CreateRole",
		Arguments: map[string]interface{}{
			"role": CreateRoleRequest{Name: "Read Only", Permissions: map[string]Permission{"*:own:list:*": Allowed, "*:own:get:*": Allowed, "*:own:*:*": Allowed}},
		},
	},
	{
		Group:        "user_roles",
		Path:         getAssignedRolesPath,
		Method:       "GET",
		FunctionName: "GetAssignedRoles",
		Arguments: map[string]interface{}{
			"userId": testUserId,
		},
	},
	{
		Group:        "user_roles",
		Path:         getAssignedRoleIDsPath,
		Method:       "GET",
		FunctionName: "GetAssignedRoleIDs",
		Arguments: map[string]interface{}{
			"userId": testUserId,
		},
	},
	{
		Group:        "user_roles",
		Path:         getUserPermissionsPath,
		Method:       "GET",
		FunctionName: "GetUserPermissions",
		Arguments: map[string]interface{}{
			"userId": testUserId,
		},
	},
}

func Test_MakeRequestError(t *testing.T) {
	for _, tt := range tests {
		setup()
		defer teardown()

		mux.HandleFunc(tt.Path, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, tt.Method, r.Method, "Expected method '%s', got %s", tt.Method, r.Method)
			w.WriteHeader(http.StatusUnauthorized)
		})

		res := callFunction(client, tt.FunctionName, tt.Arguments)
		err, _ := res[1].Interface().(error)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
	}
}
func Test_UnmarshalError(t *testing.T) {
	for _, tt := range tests {
		setup()
		defer teardown()

		mux.HandleFunc(tt.Path, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, tt.Method, r.Method, "Expected method '%s', got %s", tt.Method, r.Method)
			fmt.Fprintf(w, "not json")
		})

		res := callFunction(client, tt.FunctionName, tt.Arguments)
		err, _ := res[1].Interface().(error)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), testUnmarshalError)
	}
}
