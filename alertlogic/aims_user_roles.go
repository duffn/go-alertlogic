package alertlogic

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type RoleIdsList struct {
	RoleIds []string `json:"role_ids,omitempty"`
}

type PermissionsList struct {
	Permissions []map[string]Permission `json:"permissions,omitempty"`
}

// GetAssignedRoles gets all roles assigned to a user.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Roles_Resources-GetUserRoles
func (api *API) GetAssignedRoles(userId string) (RolesList, error) {
	return api.getRoles(fmt.Sprintf("%s/%s/users/%s/roles", aimsServicePath, api.AccountID, userId))
}

// GetAssignedRoleIDs gets the IDs for all roles assigned to a user.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Roles_Resources-GetUserRoleIds
func (api *API) GetAssignedRoleIDs(userId string) (RoleIdsList, error) {
	res, _, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/users/%s/role_ids", aimsServicePath, api.AccountID, userId), nil, nil, nil)
	if err != nil {
		return RoleIdsList{}, errors.Wrap(err, errMakeRequestError)
	}

	var r RoleIdsList
	err = json.Unmarshal(res, &r)
	if err != nil {
		return RoleIdsList{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// GetUserPermissions gets all permissions attached to a user.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Roles_Resources-GetUserPermissions
func (api *API) GetUserPermissions(userId string) (PermissionsList, error) {
	res, _, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/users/%s/permissions", aimsServicePath, api.AccountID, userId), nil, nil, nil)
	if err != nil {
		return PermissionsList{}, errors.Wrap(err, errMakeRequestError)
	}

	var r PermissionsList
	err = json.Unmarshal(res, &r)
	if err != nil {
		return PermissionsList{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// GrantUserRole grants a role to a user.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Roles_Resources-GrantUserRole
func (api *API) GrantUserRole(userId string, roleId string) (int, error) {
	return api.userRoleAssignment("PUT", userId, roleId)
}

// RevokeUserRole revokes a role from a user.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Roles_Resources-RevokeRole
func (api *API) RevokeUserRole(userId string, roleId string) (int, error) {
	return api.userRoleAssignment("DELETE", userId, roleId)
}

// userRoleAssignment has shared functionality for granting or revoking user roles.
func (api *API) userRoleAssignment(method string, userId string, roleId string) (int, error) {
	_, statusCode, err := api.makeRequest(method, fmt.Sprintf("%s/%s/users/%s/roles/%s", aimsServicePath, api.AccountID, userId, roleId), nil, nil, nil)

	if err != nil {
		return statusCode, errors.Wrap(err, errMakeRequestError)
	}

	return statusCode, nil
}
