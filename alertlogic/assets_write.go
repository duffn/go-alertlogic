package alertlogic

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

// CreateExternalDNSNameAsset creates a new asset of the type `external-dns-name` for AWS.
//
// API reference: https://console.cloudinsight.alertlogic.com/api/assets_write/#api-DeclareModify-DeclareAsset
// func (api *API) CreateExternalDNSNameAsset(deploymentId string, dnsName string) (ExternalDNSAssetResponse, error) {
// 	asset := ExternalDNSAssetRequest{
// 		Operation: "declare_asset",
// 		Type:      "external-dns-name",
// 		Scope:     "aws",
// 		Key:       fmt.Sprintf("/external-dns-name/%s", dnsName),
// 		Properties: map[string]string{
// 			"dns_name": dnsName,
// 			"name":     dnsName,
// 		},
// 	}

// 	res, _, err := api.makeRequest("POST", fmt.Sprintf("%s/%s/deployments/%s/assets", assetsWriteServicePath, api.AccountID, deploymentId), nil, nil, asset)

// 	if err != nil {
// 		return ExternalDNSAssetResponse{}, errors.Wrap(err, errMakeRequestError)
// 	}

// 	var r ExternalDNSAssetResponse
// 	err = json.Unmarshal(res, &r)
// 	if err != nil {
// 		return ExternalDNSAssetResponse{}, errors.Wrap(err, errUnmarshalError)
// 	}

// 	return r, nil
// }
