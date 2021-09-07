package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	getExternalDNSNameAssetsPath = fmt.Sprintf("/%s/%s/assets", assetsQueryServicePath, testAccountId)
)

func TestAssetsQuery_GetExternalDNSNameAssets(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"rows": 2,
		"assets": [
		[
			{
				"version": 11516,
				"type": "external-dns-name",
				"threatiness": 89.78169999999997,
				"threat_level": 3,
				"tags": {},
				"tag_keys": {},
				"state": "new",
				"scope_external_scan_request_id": "deebd20e-9814-4256-9738-7cfa757c94ae",
				"scope_external_last_external_scan_time": 1631034658,
				"scope_external_last_dequeue_time": 1631029065,
				"scope_external_heartbeat": 1631034599,
				"scope_aws_state": "new",
				"scope_aws_name": "9876-qwert.elb.us-east-1.amazonaws.com",
				"scope_aws_dns_name": "9876-qwert.elb.us-east-1.amazonaws.com",
				"native_type": "external-dns-name",
				"name": "9876-qwert.elb.us-east-1.amazonaws.com",
				"modified_on": 1631034658921,
				"key": "/external-dns-name/9876-qwert.elb.us-east-1.amazonaws.com",
				"dns_name": "9876-qwert.elb.us-east-1.amazonaws.com",
				"deployment_id": "f69f7395-ce6a-43af-b641-e4bad8bbec88",
				"deleted_on": 0,
				"declared": true,
				"created_on": 1620924080257,
				"account_id": "12345678"
			}
		],
		[
			{
				"version": 9835,
				"type": "external-dns-name",
				"threatiness": 93.13829999999996,
				"threat_level": 3,
				"tags": {},
				"tag_keys": {},
				"state": "new",
				"scope_external_scan_request_id": "d0b4844d-47f6-4305-9723-c49c8e836252",
				"scope_external_last_external_scan_time": 1630996539,
				"scope_external_last_dequeue_time": 1630991137,
				"scope_external_heartbeat": 1630996480,
				"scope_aws_state": "new",
				"scope_aws_name": "abcd-1234.elb.us-east-1.amazonaws.com",
				"scope_aws_dns_name": "abcd-1234.elb.us-east-1.amazonaws.com",
				"native_type": "external-dns-name",
				"name": "abcd-1234.elb.us-east-1.amazonaws.com",
				"modified_on": 1630996539988,
				"key": "/external-dns-name/abcd-1234.elb.us-east-1.amazonaws.com",
				"dns_name": "abcd-1234.elb.us-east-1.amazonaws.com",
				"deployment_id": "11ef5a1e-74ac-4bc3-a082-1fa8b64abb64",
				"deleted_on": 0,
				"declared": true,
				"created_on": 1620910786980,
				"account_id": "12345678"
			}
		]
	]}`

	mux.HandleFunc(getExternalDNSNameAssetsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := ExternalDNSNameAssets{
		Rows: 2,
		ExternalDNSAssets: [][]ExternalDNSNameAsset{
			{
				{
					Version:                           11516,
					Type:                              "external-dns-name",
					Threatiness:                       89.78169999999997,
					ThreatLevel:                       3,
					Tags:                              Tag{},
					TagKeys:                           Tag{},
					State:                             "new",
					ScopeExternalScanRequestID:        "deebd20e-9814-4256-9738-7cfa757c94ae",
					ScopeExternalLastExternalScanTime: 1631034658,
					ScopeExternalLastDequeueTime:      1631029065,
					ScopeExternalHeartbeat:            1631034599,
					ScopeAwsState:                     "new",
					ScopeAwsName:                      "9876-qwert.elb.us-east-1.amazonaws.com",
					ScopeAwsDNSName:                   "9876-qwert.elb.us-east-1.amazonaws.com",
					NativeType:                        "external-dns-name",
					Name:                              "9876-qwert.elb.us-east-1.amazonaws.com",
					ModifiedOn:                        1631034658921,
					Key:                               "/external-dns-name/9876-qwert.elb.us-east-1.amazonaws.com",
					DNSName:                           "9876-qwert.elb.us-east-1.amazonaws.com",
					DeploymentID:                      "f69f7395-ce6a-43af-b641-e4bad8bbec88",
					DeletedOn:                         0,
					Declared:                          true,
					CreatedOn:                         1620924080257,
					AccountID:                         "12345678",
				},
			},
			{
				{
					Version:                           9835,
					Type:                              "external-dns-name",
					Threatiness:                       93.13829999999996,
					ThreatLevel:                       3,
					Tags:                              Tag{},
					TagKeys:                           Tag{},
					State:                             "new",
					ScopeExternalScanRequestID:        "d0b4844d-47f6-4305-9723-c49c8e836252",
					ScopeExternalLastExternalScanTime: 1630996539,
					ScopeExternalLastDequeueTime:      1630991137,
					ScopeExternalHeartbeat:            1630996480,
					ScopeAwsState:                     "new",
					ScopeAwsName:                      "abcd-1234.elb.us-east-1.amazonaws.com",
					ScopeAwsDNSName:                   "abcd-1234.elb.us-east-1.amazonaws.com",
					NativeType:                        "external-dns-name",
					Name:                              "abcd-1234.elb.us-east-1.amazonaws.com",
					ModifiedOn:                        1630996539988,
					Key:                               "/external-dns-name/abcd-1234.elb.us-east-1.amazonaws.com",
					DNSName:                           "abcd-1234.elb.us-east-1.amazonaws.com",
					DeploymentID:                      "11ef5a1e-74ac-4bc3-a082-1fa8b64abb64",
					DeletedOn:                         0,
					Declared:                          true,
					CreatedOn:                         1620910786980,
					AccountID:                         "12345678",
				},
			},
		},
	}

	assets, err := client.GetExternalDNSNameAssets()
	if assert.NoError(t, err) {
		assert.Equal(t, assets, want)
	}
}

func TestAssetsQuery_GetExternalDNSNameAssetsMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getExternalDNSNameAssetsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetExternalDNSNameAssets()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestAssetsQuery_GetExternalDNSNameAssetsUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getExternalDNSNameAssetsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetExternalDNSNameAssets()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}
