package alertlogic

import (
	"fmt"

	"github.com/pkg/errors"
)

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

// CreateExternalDNSNameAsset creates a new asset of the type `external-dns-name` for AWS.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/assets_write/#api-DeclareModify-DeclareAsset
func (api *API) CreateExternalDNSNameAsset(deploymentId string, dnsName string) (int, error) {
	return api.modifyExternalDNSNameAsset(deploymentId, dnsName, "")
}

// UpdateExternalDNSNameAsset updates an existing asset of the type `external-dns-name` for AWS.
// The only item that you may update is the DNS name of the asset.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/assets_write/#api-DeclareModify-DeclareAsset
func (api *API) UpdateExternalDNSNameAsset(deploymentId string, dnsName string, oldDnsName string) (int, error) {
	return api.modifyExternalDNSNameAsset(deploymentId, dnsName, oldDnsName)
}

// RemoveExternalDNSNameAsset creates a new asset of the type `external-dns-name` for AWS.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/assets_write/#api-DeclareModify-RemoveAsset
func (api *API) RemoveExternalDNSNameAsset(deploymentId string, dnsName string) (int, error) {
	asset := ExternalDNSAssetRequest{
		Operation: "remove_asset",
		Type:      "external-dns-name",
		Scope:     "aws",
		Key:       fmt.Sprintf("/external-dns-name/%s", dnsName),
	}

	_, statusCode, err := api.makeRequest("PUT", fmt.Sprintf("%s/%s/deployments/%s/assets", assetsWriteServicePath, api.AccountID, deploymentId), nil, nil, asset)

	if err != nil {
		return statusCode, errors.Wrap(err, errMakeRequestError)
	}

	return statusCode, nil
}

// modifyExternalDNSNameAsset holds shared logic for creating or modifying an external DNS
// asset.
func (api *API) modifyExternalDNSNameAsset(deploymentId string, dnsName string, oldDnsName string) (int, error) {
	keyDns := dnsName
	if oldDnsName != "" {
		keyDns = oldDnsName
	}

	asset := ExternalDNSAssetRequest{
		Operation: "declare_asset",
		Type:      "external-dns-name",
		Scope:     "aws",
		Key:       fmt.Sprintf("/external-dns-name/%s", keyDns),
		Properties: map[string]string{
			"dns_name": dnsName,
			"name":     dnsName,
			"state":    "new",
		},
	}

	_, statusCode, err := api.makeRequest("PUT", fmt.Sprintf("%s/%s/deployments/%s/assets", assetsWriteServicePath, api.AccountID, deploymentId), nil, nil, asset)

	if err != nil {
		return statusCode, errors.Wrap(err, errMakeRequestError)
	}

	return statusCode, nil
}
