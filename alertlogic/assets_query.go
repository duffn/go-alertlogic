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

type ExternalDNSAssetRequest struct {
	Operation     string            `json:"operation,omitempty"`
	Type          string            `json:"type,omitempty"`
	Scope         string            `json:"scope,omitempty"`
	Properties    map[string]string `json:"properties,omitempty"`
	Relationships []Relationship    `json:"relationships,omitempty"`
	Key           string            `json:"key,omitempty"`
}

type Relationship struct {
	Key  string `json:"key,omitempty"`
	Type string `json:"type,omitempty"`
}

type ExternalDNSAssetResponse struct {
	AccountID     string         `json:"account_id,omitempty"`
	DeploymentID  string         `json:"deployment_id,omitempty"`
	Key           string         `json:"key,omitempty"`
	Relationships []Relationship `json:"relationships,omitempty"`
	Type          string         `json:"type,omitempty"`
}

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

// CreateExternalDNSNameAsset creates a new asset of the type `external-dns-name` for AWS.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/assets_write/#api-DeclareModify-DeclareAsset
func (api *API) CreateExternalDNSNameAsset(deploymentId string, dnsName string) (ExternalDNSAssetResponse, error) {
	asset := ExternalDNSAssetRequest{
		Operation: "declare_asset",
		Type:      "external-dns-name",
		Scope:     "aws",
		Key:       fmt.Sprintf("/external-dns-name/%s", dnsName),
		Properties: map[string]string{
			"dns_name": dnsName,
			"name":     dnsName,
		},
	}

	res, _, err := api.makeRequest("POST", fmt.Sprintf("%s/%s/deployments/%s/assets", assetsWriteServicePath, api.AccountID, deploymentId), nil, nil, asset)

	if err != nil {
		return ExternalDNSAssetResponse{}, errors.Wrap(err, errMakeRequestError)
	}

	var r ExternalDNSAssetResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return ExternalDNSAssetResponse{}, errors.Wrap(err, errUnmarshalError)
	}

	return r, nil
}
