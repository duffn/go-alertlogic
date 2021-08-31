package alertlogic

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

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

// GetAccountDetails gets details of an account.
func (api *API) GetAccountDetails() (AccountDetailsResponse, error) {
	res, _, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/account", aimsServicePath, api.AccountID), nil, nil, nil)

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
	_, statusCode, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/accounts/%s/%s", aimsServicePath, api.AccountID, accountRelationship, relatedAccountId), nil, nil, nil)

	if err != nil {
		return statusCode, errors.Wrap(err, errMakeRequestError)
	}

	return statusCode, nil
}
