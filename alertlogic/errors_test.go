package alertlogic

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func callFunction(obj interface{}, fn string, args []interface{}) (res []reflect.Value) {
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
	Arguments    []interface{}
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
		Arguments: []interface{}{
			UpdateAccountDetailsRequest{MfaRequired: false},
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
		Arguments: []interface{}{
			testRoleId,
		},
	},
	{
		Group:        "roles",
		Path:         getGlobalRoleDetailsPath,
		Method:       "GET",
		FunctionName: "GetGlobalRoleDetails",
		Arguments: []interface{}{
			testRoleId,
		},
	},
	{
		Group:        "user_roles",
		Path:         getAssignedRolesPath,
		Method:       "GET",
		FunctionName: "GetAssignedRoles",
		Arguments: []interface{}{
			testUserId,
		},
	},
	{
		Group:        "user_roles",
		Path:         getAssignedRoleIDsPath,
		Method:       "GET",
		FunctionName: "GetAssignedRoleIDs",
		Arguments: []interface{}{
			testUserId,
		},
	},
	{
		Group:        "user_roles",
		Path:         getUserPermissionsPath,
		Method:       "GET",
		FunctionName: "GetUserPermissions",
		Arguments: []interface{}{
			testUserId,
		},
	},
	{
		Group:        "users",
		Path:         authenticatePath,
		Method:       "POST",
		FunctionName: "Authenticate",
		Arguments:    nil,
	},
	{
		Group:        "users",
		Path:         createUserPath,
		Method:       "POST",
		FunctionName: "CreateUser",
		Arguments: []interface{}{
			CreateUserRequest{Email: testEmail, Name: "Bob Loblaw"},
			false,
		},
	},
	{
		Group:        "users",
		Path:         listUsersByEmailPath,
		Method:       "GET",
		FunctionName: "ListUsersByEmail",
		Arguments: []interface{}{
			testEmail,
			true,
			true,
			true,
		},
	},
	{
		Group:        "users",
		Path:         getUserDetailsByIdPath,
		Method:       "GET",
		FunctionName: "GetUserDetailsById",
		Arguments: []interface{}{
			testUserId,
			true,
			true,
			true,
		},
	},
	{
		Group:        "users",
		Path:         listUsersPath,
		Method:       "GET",
		FunctionName: "ListUsers",
		Arguments: []interface{}{
			true,
			true,
			true,
			"",
		},
	},
	{
		Group:        "users",
		Path:         updateUserPath,
		Method:       "POST",
		FunctionName: "UpdateUserDetails",
		Arguments: []interface{}{
			testUserId,
			UpdateUserRequest{Email: "new@email.com"},
			false,
		},
	},
	{
		Group:        "users",
		Path:         getUserDetailsByUsernamePath,
		Method:       "GET",
		FunctionName: "GetUserDetailsByUsername",
		Arguments: []interface{}{
			testUserId,
			true,
			true,
			true,
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
