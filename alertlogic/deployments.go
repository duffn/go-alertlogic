package alertlogic

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Deployment struct {
	Version       int64           `json:"version,omitempty"`
	Status        Status          `json:"status,omitempty"`
	Scope         Scope           `json:"scope,omitempty"`
	Scan          bool            `json:"scan,omitempty"`
	Platform      Platform        `json:"platform,omitempty"`
	Name          string          `json:"name,omitempty"`
	Modified      ModifiedCreated `json:"modified,omitempty"`
	Mode          string          `json:"mode,omitempty"`
	ID            string          `json:"id,omitempty"`
	Enabled       bool            `json:"enabled,omitempty"`
	Discover      bool            `json:"discover,omitempty"`
	Credentials   []Credential    `json:"credentials,omitempty"`
	Created       ModifiedCreated `json:"created,omitempty"`
	CloudDefender CloudDefender   `json:"cloud_defender,omitempty"`
	AccountID     string          `json:"account_id,omitempty"`
}

type CloudDefender struct {
	Enabled    bool   `json:"enabled,omitempty"`
	LocationID string `json:"location_id,omitempty"`
}

type Credential struct {
	ID      string `json:"id,omitempty"`
	Purpose string `json:"purpose,omitempty"`
	Version string `json:"version,omitempty"`
}

type Platform struct {
	Type    string  `json:"type,omitempty"`
	ID      string  `json:"id,omitempty"`
	Monitor Monitor `json:"monitor,omitempty"`
	Default bool    `json:"default,omitempty"`
}

type Monitor struct {
	Enabled         bool   `json:"enabled,omitempty"`
	CTInstallRegion string `json:"ct_install_region,omitempty"`
}

type Scope struct {
	Include []Include `json:"include,omitempty"`
	Exclude []Exclude `json:"exclude,omitempty"`
}

type Exclude struct {
	Type string `json:"type,omitempty"`
	Key  string `json:"key,omitempty"`
}

type Include struct {
	Type   string `json:"type,omitempty"`
	Key    string `json:"key,omitempty"`
	Policy Policy `json:"policy,omitempty"`
}

type Policy struct {
	ID string `json:"id,omitempty"`
}

type Status struct {
	Status  string `json:"status,omitempty"`
	Updated int64  `json:"updated,omitempty"`
}

// ListDeployments lists deployments for an account.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/deployments/#api-Resources-ListDeployments
func (api *API) ListDeployments() ([]Deployment, error) {
	res, _, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/deployments", deploymentServicePath, api.AccountID), nil, nil, nil)
	if err != nil {
		return []Deployment{}, errors.Wrap(err, errMakeRequestError)
	}

	var r []Deployment
	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Deployment{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}

// GetDeployment gets a single deployment for an account.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/deployments/#api-Resources-GetDeployment
func (api *API) GetDeployment(deploymentId string) (Deployment, error) {
	res, _, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/deployments/%s", deploymentServicePath, api.AccountID, deploymentId), nil, nil, nil)
	if err != nil {
		return Deployment{}, errors.Wrap(err, errMakeRequestError)
	}

	var r Deployment
	err = json.Unmarshal(res, &r)
	if err != nil {
		return Deployment{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}
