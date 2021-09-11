package alertlogic

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type ExternalDNSNameAssets struct {
	Rows              int64                    `json:"rows"`
	ExternalDNSAssets [][]ExternalDNSNameAsset `json:"assets"`
}

type ExternalDNSNameAsset struct {
	Version                           int64   `json:"version,omitempty"`
	Type                              string  `json:"type,omitempty"`
	Threatiness                       float64 `json:"threatiness,omitempty"`
	ThreatLevel                       int64   `json:"threat_level,omitempty"`
	Tags                              Tag     `json:"tags,omitempty"`
	TagKeys                           Tag     `json:"tag_keys,omitempty"`
	State                             string  `json:"state,omitempty"`
	ScopeExternalScanRequestID        string  `json:"scope_external_scan_request_id,omitempty"`
	ScopeExternalLastExternalScanTime int64   `json:"scope_external_last_external_scan_time,omitempty"`
	ScopeExternalLastDequeueTime      int64   `json:"scope_external_last_dequeue_time,omitempty"`
	ScopeExternalHeartbeat            int64   `json:"scope_external_heartbeat,omitempty"`
	ScopeAwsState                     string  `json:"scope_aws_state,omitempty"`
	ScopeAwsName                      string  `json:"scope_aws_name,omitempty"`
	ScopeAwsDNSName                   string  `json:"scope_aws_dns_name,omitempty"`
	NativeType                        string  `json:"native_type,omitempty"`
	Name                              string  `json:"name,omitempty"`
	ModifiedOn                        int64   `json:"modified_on,omitempty"`
	Key                               string  `json:"key,omitempty"`
	DNSName                           string  `json:"dns_name,omitempty"`
	DeploymentID                      string  `json:"deployment_id,omitempty"`
	DeletedOn                         int64   `json:"deleted_on,omitempty"`
	Declared                          bool    `json:"declared,omitempty"`
	CreatedOn                         int64   `json:"created_on,omitempty"`
	AccountID                         string  `json:"account_id,omitempty"`
}

type Tag struct{}

// GetExternalDNSNamesAssets gets external DNS assets for an account.
// This endpoint only gets assets that are of the `external-dns-name` type.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/assets_query/#api-Queries-QueryAccountAssets
func (api *API) GetExternalDNSNameAssets() (ExternalDNSNameAssets, error) {
	params := map[string]string{"asset_types": "e:external-dns-name"}
	res, _, err := api.makeRequest("GET", fmt.Sprintf("%s/%s/assets", assetsQueryServicePath, api.AccountID), nil, params, nil)

	if err != nil {
		return ExternalDNSNameAssets{}, errors.Wrap(err, errMakeRequestError)
	}

	var r ExternalDNSNameAssets
	err = json.Unmarshal(res, &r)
	if err != nil {
		return ExternalDNSNameAssets{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}
