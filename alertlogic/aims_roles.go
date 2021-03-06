package alertlogic

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type RolesList struct {
	Roles []Role `json:"roles,omitempty"`
}

type Role struct {
	ID                string                `json:"id,omitempty"`
	AccountID         string                `json:"account_id,omitempty"`
	Name              string                `json:"name"`
	Permissions       map[string]Permission `json:"permissions"`
	Version           int64                 `json:"version,omitempty"`
	Global            bool                  `json:"global,omitempty"`
	LegacyPermissions []string              `json:"legacy_permissions,omitempty"`
	Created           ModifiedCreated       `json:"created,omitempty"`
	Modified          ModifiedCreated       `json:"modified,omitempty"`
}

type Permission string

const (
	Allowed Permission = "allowed"
	Denied  Permission = "denied"
)

// GetRoleDetails retrieves a role's details.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Role_Resources-GetRole
func (api *API) GetRoleDetails(roleId string) (Role, error) {
	return api.getRole(fmt.Sprintf("%s/%s/roles/%s", aimsServicePath, api.AccountID, roleId))
}

// GetRoleDetails retrieves a global role's details.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Role_Resources-GetGlobalRole
func (api *API) GetGlobalRoleDetails(roleId string) (Role, error) {
	return api.getRole(fmt.Sprintf("%s/roles/%s", aimsServicePath, roleId))
}

// ListRoles list all roles for a specific account.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Role_Resources-ListRoles
func (api *API) ListRoles() (RolesList, error) {
	return api.getRoles(fmt.Sprintf("%s/%s/roles", aimsServicePath, api.AccountID))
}

// ListGlobalRoles list all roles for across all accounts.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Role_Resources-ListGlobalRoles
func (api *API) ListGlobalRoles() (RolesList, error) {
	return api.getRoles(fmt.Sprintf("%s/roles", aimsServicePath))
}

// getRole holds shared logic for retrieving a Rols from the API.
func (api *API) getRole(path string) (Role, error) {
	res, _, err := api.makeRequest("GET", path, nil, nil, nil)
	if err != nil {
		return Role{}, errors.Wrap(err, errMakeRequestError)
	}

	var r Role
	err = json.Unmarshal(res, &r)
	if err != nil {
		return Role{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// getRoles holds shared logic for retrieving multiple Roles from the API.
func (api *API) getRoles(path string) (RolesList, error) {
	res, _, err := api.makeRequest("GET", path, nil, nil, nil)
	if err != nil {
		return RolesList{}, errors.Wrap(err, errMakeRequestError)
	}

	var r RolesList
	err = json.Unmarshal(res, &r)
	if err != nil {
		return RolesList{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}
