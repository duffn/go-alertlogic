package alertlogic

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// aimsServicePath is the path for the aims service.
const aimsServicePath = "aims/v1"

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

// Account is the account level information.
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
	ID          string          `json:"id,omitempty"`
	AccountID   string          `json:"account_id,omitempty"`
	Name        string          `json:"name,omitempty"`
	Username    string          `json:"username,omitempty"`
	Email       string          `json:"email,omitempty"`
	Active      bool            `json:"active,omitempty"`
	Locked      bool            `json:"locked,omitempty"`
	Version     int64           `json:"version,omitempty"`
	LinkedUsers []LinkedUser    `json:"linked_users,omitempty"`
	Created     ModifiedCreated `json:"created,omitempty"`
	Modified    ModifiedCreated `json:"modified,omitempty"`
}

// LinkedUser are any users linked to the current user.
type LinkedUser struct {
	UserID   int64  `json:"user_id,omitempty"`
	Location string `json:"location,omitempty"`
}

// AccountDetailsResponse represents account details.
type AccountDetailsResponse struct {
	ID                  string          `json:"id,omitempty"`
	Name                string          `json:"name,omitempty"`
	Active              bool            `json:"active,omitempty"`
	Version             int64           `json:"version,omitempty"`
	AccessibleLocations []string        `json:"accessible_locations,omitempty"`
	DefaultLocation     string          `json:"default_location,omitempty"`
	Created             ModifiedCreated `json:"created,omitempty"`
	Modified            ModifiedCreated `json:"modified,omitempty"`
}

// AccountRelationship is the relationship of one account to another.
type AccountRelationship string

const (
	BillsTo  AccountRelationship = "bills_to"
	Managed  AccountRelationship = "managed"
	Managing AccountRelationship = "managing"
)

// Authenticate authenticates a user and returns a token and user details. If you're using
// this method directly, then returned token should be used as API.APIToken for all future
// calls to the API.
// Preferably, you should use `NewWithUsernameAndPassword` which will authenticate with the
// API and set your token on API for future calls.
func (api *API) Authenticate() (AuthenticateResponse, error) {
	res, _, err := api.makeRequest("POST", fmt.Sprintf("%s/authenticate", aimsServicePath), nil, nil)

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

// GetAccountDetails gets details of an account.
func (api *API) GetAccountDetails() (AccountDetailsResponse, error) {
	res, _, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/account", aimsServicePath, api.AccountID), nil, nil)

	if err != nil {
		return AccountDetailsResponse{}, errors.Wrap(err, errMakeRequestError)
	}

	var r AccountDetailsResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return AccountDetailsResponse{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// GetAccountRelationship gets a relationship between one account and another. The first account, provided
// when the API client is created, is the primary account, and `relatedAccountId` is the secondary account.
// This API returns 204 when these two accounts have an `accountRelationship` relationship and 404 when
// they do not.
func (api *API) GetAccountRelationship(relatedAccountId string, accountRelationship AccountRelationship) (int, error) {
	_, statusCode, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/accounts/%s/%s", aimsServicePath, api.AccountID, accountRelationship, relatedAccountId), nil, nil)

	if err != nil {
		return statusCode, errors.Wrap(err, errMakeRequestError)
	}

	return statusCode, nil
}
