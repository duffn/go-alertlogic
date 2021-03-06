package alertlogic

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// AuthenticateResponse holds the response from Authenticate.
type AuthenticateResponse struct {
	Authentication Authentication `json:"authentication"`
}

// Authentication is the authentication information.
type Authentication struct {
	User            User    `json:"user"`
	Account         Account `json:"account"`
	Token           string  `json:"token"`
	TokenExpiration int64   `json:"token_expiration"`
}

// Account is the account level information for a user.
type Account struct {
	ID                  string          `json:"id,omitempty"`
	Name                string          `json:"name,omitempty"`
	Active              bool            `json:"active,omitempty"`
	Version             int64           `json:"version,omitempty"`
	AccessibleLocations []string        `json:"accessible_locations,omitempty"`
	DefaultLocation     string          `json:"default_location,omitempty"`
	DefaultESBLocation  string          `json:"default_esb_location,omitempty"`
	Created             ModifiedCreated `json:"created,omitempty"`
	Modified            ModifiedCreated `json:"modified,omitempty"`
}

// User is the user level information.
type User struct {
	ID             string          `json:"id,omitempty"`
	AccountID      string          `json:"account_id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Username       string          `json:"username,omitempty"`
	Email          string          `json:"email,omitempty"`
	Active         bool            `json:"active,omitempty"`
	Locked         bool            `json:"locked,omitempty"`
	Version        int64           `json:"version,omitempty"`
	MfaEnabled     *bool           `json:"mfa_enabled,omitempty"`
	MobilePhone    *string         `json:"mobile_phone,omitempty"`
	LinkedUsers    []LinkedUser    `json:"linked_users,omitempty"`
	UserCredential *UserCredential `json:"user_credential,omitempty"`
	RoleIds        *[]string       `json:"role_ids,omitempty"`
	AccessKeys     *[]AccessKey    `json:"access_keys,omitempty"`
	Created        ModifiedCreated `json:"created,omitempty"`
	Modified       ModifiedCreated `json:"modified,omitempty"`
}

// UserCredential is a user's credential information.
type UserCredential struct {
	Version         int             `json:"version,omitempty"`
	OneTimePassword bool            `json:"one_time_password,omitempty"`
	LastLogin       int             `json:"last_login,omitempty"`
	Created         ModifiedCreated `json:"created,omitempty"`
	Modified        ModifiedCreated `json:"modified,omitempty"`
}

type AccessKey struct {
	Label       string          `json:"label,omitempty"`
	LastLogin   int             `json:"last_login,omitempty"`
	Created     ModifiedCreated `json:"created,omitempty"`
	Modified    ModifiedCreated `json:"modified,omitempty"`
	AccessKeyID string          `json:"access_key_id,omitempty"`
}

// LinkedUser are any users linked to the current user.
type LinkedUser struct {
	UserID   int64  `json:"user_id,omitempty"`
	Location string `json:"location,omitempty"`
}

// CreateUserRequest holds the user create request data.
type CreateUserRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	RoleId            string `json:"role_id,omitempty"`
	Active            bool   `json:"active,omitempty"`
	MobilePhone       string `json:"mobile_phone,omitempty"`
	Phone             string `json:"phone,omitempty"`
	WebhookUrl        string `json:"webhook_url,omitempty"`
	NotificationsOnly bool   `json:"notifications_only,omitempty"`
}

// UpdateUserRequest holds the user update request data.
type UpdateUserRequest = CreateUserRequest

// ListUsersByEmailResponse holds the response from list users by email.
type UserList struct {
	Users []User `json:"users"`
}

// Authenticate authenticates a user and returns a token and user details. If you're using
// this method directly, then returned token should be used as API.APIToken for all future
// calls to the API.
// Preferably, you should use `NewWithUsernameAndPassword` which will authenticate with the
// API and set your token on API for future calls.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Authentication_and_Authorization_Resources-Authenticate
func (api *API) Authenticate() (AuthenticateResponse, error) {
	res, _, err := api.makeRequest("POST", fmt.Sprintf("%s/authenticate", aimsServicePath), nil, nil, nil)

	if err != nil {
		return AuthenticateResponse{}, errors.Wrap(err, errMakeRequestError)
	}

	var r AuthenticateResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return AuthenticateResponse{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// CreateUser creates a new user.
// If true, `oneTimePassword` will set the user's password as a one-time password and require them
// to supply a new password upon first login.
// If `Password` is not supplied in the `CreateUserRequest`, the user will be emailed a link to
// set their password.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-CreateUser
func (api *API) CreateUser(user CreateUserRequest, oneTimePassword bool) (User, error) {
	if oneTimePassword && user.Password == "" {
		return User{}, errors.New("oneTimePassword must be accompanied by CreateUserRequest.Password")
	}

	var params map[string]string
	if oneTimePassword {
		params = map[string]string{"one_time_password": "true"}
	}

	res, _, err := api.makeRequest("POST", fmt.Sprintf("%s/%s/users", aimsServicePath, api.AccountID), nil, params, user)

	if err != nil {
		return User{}, errors.Wrap(err, errMakeRequestError)
	}

	var r User
	err = json.Unmarshal(res, &r)
	if err != nil {
		return User{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// DeleteUser deletes a user.
// Note that this endpoint returns a 204 status code even if the user ID does not exist and only
// returns a 400 error if you try to delete the account associated with your API token.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-DeleteUser
func (api *API) DeleteUser(userId string) (int, error) {
	_, statusCode, err := api.makeRequest("DELETE", fmt.Sprintf("%s/%s/users/%s", aimsServicePath, api.AccountID, userId), nil, nil, nil)

	if err != nil {
		return statusCode, errors.Wrap(err, errMakeRequestError)
	}

	return statusCode, nil
}

// ListUsersByEmail retrieves users by email address.
// You can include access keys, credentials, or role IDs with includeAccessKeys,
// includeUserCredentials, and includeRoleIds respectively.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-ListUsersByEmail
func (api *API) ListUsersByEmail(email string, includeAccessKeys bool, includeUserCredentials bool, includeRoleIds bool) (UserList, error) {
	return api.getUsers(fmt.Sprintf("%s/users/email/%s", aimsServicePath, url.QueryEscape(email)), includeAccessKeys, includeUserCredentials, includeRoleIds, "")
}

// GetUserDetailsById retrieves a user's details by their ID.
// You can include access keys, credentials, or role IDs with includeAccessKeys,
// includeUserCredentials, and includeRoleIds respectively.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-GetUserDetailsByUserId
func (api *API) GetUserDetailsById(userId string, includeAccessKeys bool, includeUserCredentials bool, includeRoleIds bool) (User, error) {
	return api.getUser(fmt.Sprintf("%s/user/%s", aimsServicePath, userId), includeAccessKeys, includeUserCredentials, includeRoleIds)
}

// GetUserDetails retrieves a user's details by their ID.
// You can include access keys, credentials, or role IDs with includeAccessKeys,
// includeUserCredentials, and includeRoleIds respectively.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-GetUserDetails
func (api *API) GetUserDetails(userId string, includeAccessKeys bool, includeUserCredentials bool, includeRoleIds bool) (User, error) {
	return api.getUser(fmt.Sprintf("%s/%s/users/%s", aimsServicePath, api.AccountID, userId), includeAccessKeys, includeUserCredentials, includeRoleIds)
}

// ListUsersByEmail retrieves users by email address.
// You can include access keys, credentials, or role IDs with includeAccessKeys,
// includeUserCredentials, and includeRoleIds respectively.
// Specifying roleId will only return users that belong to that role.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-ListUsers
func (api *API) ListUsers(includeAccessKeys bool, includeUserCredentials bool, includeRoleIds bool, roleId string) (UserList, error) {
	return api.getUsers(fmt.Sprintf("%s/%s/users", aimsServicePath, api.AccountID), includeAccessKeys, includeUserCredentials, includeRoleIds, roleId)
}

// UpdateUserDetails updates a user.
// If true, `oneTimePassword` will set the user's password as a one-time password and require them
// to supply a new password upon first login.
// If `Password` is not supplied in the `UpdateUserRequest`, the user will be emailed a link to
// set their password.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-UpdateUser
func (api *API) UpdateUserDetails(userId string, user UpdateUserRequest, oneTimePassword bool) (User, error) {
	if oneTimePassword && user.Password == "" {
		return User{}, errors.New("oneTimePassword must be accompanied by CreateUserRequest.Password")
	}

	var params map[string]string
	if oneTimePassword {
		params = map[string]string{"one_time_password": "true"}
	}

	res, _, err := api.makeRequest("POST", fmt.Sprintf("%s/%s/users/%s", aimsServicePath, api.AccountID, userId), nil, params, user)

	if err != nil {
		return User{}, errors.Wrap(err, errMakeRequestError)
	}

	var r User
	err = json.Unmarshal(res, &r)
	if err != nil {
		return User{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// GetUserDetailsByUsername retrieves a user's details by their username.
// You can include access keys, credentials, or role IDs with includeAccessKeys,
// includeUserCredentials, and includeRoleIds respectively.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-GetUserDetailsByUserId
func (api *API) GetUserDetailsByUsername(username string, includeAccessKeys bool, includeUserCredentials bool, includeRoleIds bool) (User, error) {
	return api.getUser(fmt.Sprintf("%s/user/username/%s", aimsServicePath, username), includeAccessKeys, includeUserCredentials, includeRoleIds)
}

// getUser holds shared logic for retrieving a User from the API.
func (api *API) getUser(path string, includeAccessKeys bool, includeUserCredentials bool, includeRoleIds bool) (User, error) {
	var params = map[string]string{
		"include_access_keys":     "false",
		"include_user_credential": "false",
		"include_role_ids":        "false",
	}

	if includeAccessKeys {
		params["include_access_keys"] = "true"
	}
	if includeUserCredentials {
		params["include_user_credential"] = "true"
	}
	if includeRoleIds {
		params["include_role_ids"] = "true"
	}

	res, _, err := api.makeRequest("GET", path, nil, params, nil)
	if err != nil {
		return User{}, errors.Wrap(err, errMakeRequestError)
	}

	var r User
	err = json.Unmarshal(res, &r)
	if err != nil {
		return User{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// getUsers holds shared logic for retrieving multiple Users from the API.
func (api *API) getUsers(path string, includeAccessKeys bool, includeUserCredentials bool, includeRoleIds bool, roleId string) (UserList, error) {
	var params = map[string]string{
		"include_access_keys":     "false",
		"include_user_credential": "false",
		"include_role_ids":        "false",
	}

	if includeAccessKeys {
		params["include_access_keys"] = "true"
	}
	if includeUserCredentials {
		params["include_user_credential"] = "true"
	}
	if includeRoleIds {
		params["include_role_ids"] = "true"
	}
	if roleId != "" {
		params["role_id"] = roleId
	}

	res, _, err := api.makeRequest("GET", path, nil, params, nil)
	if err != nil {
		return UserList{}, errors.Wrap(err, errMakeRequestError)
	}

	var r UserList
	err = json.Unmarshal(res, &r)
	if err != nil {
		return UserList{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}
