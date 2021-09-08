package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDeploymentId = "50668317-feb8-49d1-b401-7219bfa22417"
)

var (
	listDeploymentsPath = fmt.Sprintf("/%s/%s/deployments", deploymentServicePath, testAccountId)
	getDeploymentPath   = fmt.Sprintf("/%s/%s/deployments/%s", deploymentServicePath, testAccountId, testDeploymentId)
)

func TestDeployments_ListDeployments(t *testing.T) {
	setup()
	defer teardown()

	const response = `[
	{
		"id": "AF50CF2D-3E77-46E3-B003-0262D77B2A65",
		"account_id": "01000001",
		"name": "AWS Production Deployment",
		"platform": {
			"type": "aws",
			"id": "111111111111",
			"monitor": {
				"enabled": true,
				"ct_install_region": "us-east-1"
			}
		},
		"mode": "automatic",
		"enabled": true,
		"discover": true,
		"scan": true,
		"scope": {
			"include": [
				{
					"type": "region",
					"key": "/aws/us-east-1"
				},
				{
					"type": "region",
					"key": "/aws/us-east-2",
					"policy": {
						"id": "A8E8B104-8F45-411D-A240-A30EA5FE25B0"
					}
				},
				{
					"type": "vpc",
					"key": "/aws/us-west-1/vpc/vpc-12345678",
					"policy": {
						"id": "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"
					}	
				}
			],
			"exclude": [
				{
					"type": "region",
					"key": "/aws/ap-southeast-1"
				}
			]
		},
		"cloud_defender": {
			"enabled": true,
			"location_id": "defender-us-denver"
		},
		"credentials": [
			{
				"id": "E09F0AF8-18F8-49CE-B9AC-01C3E214B4EB",
				"purpose": "discover",
				"version": "2018-01-01"
			},
			{
				"id": "D9B5066E-6145-4D41-9B3A-C123BDE64ABD",
				"purpose": "x-account-monitor",
				"version": "2018-01-01"
			}
		],
		"status": {
			"status": "ok",
			"updated": 1493671342
		},
		"created": {
			"at": 1496673172,
			"by": "BBCAE827-22A6-433D-884E-22AABF2DC82B"
		},
		"modified": {
			"at": 1496673245,
			"by": "1443C74F-2AF7-4D7F-BF4A-FDFFE9C67182"
		},
		"version": 2
	},
	{
		"id": "86975343-2DA7-4E77-9F52-8488C4217191",
		"account_id": "01000001",
		"name": "Azure Production Deployment",
		"platform": {
			"type": "azure",
			"id": "93B0107D-21AD-48B0-9101-C60F20E09966"
		},
		"mode": "manual",
		"enabled": true,
		"discover": true,
		"scan": true,
		"scope": {
			"include": [
				{
					"type": "location",
					"key": "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/locations/centralus"
				},
				{
					"type": "location",
					"key": "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/locations/canadacentral",
					"policy": {
						"id": "A8E8B104-8F45-411D-A240-A30EA5FE25B0"
					}
				},
				{
					"type": "vnet",
					"key": "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/resourceGroups/test-rg-01/providers/Microsoft.Network/virtualNetworks/test-vnet-01",
					"policy": {
						"id": "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"
					}
				}
			],
			"exclude": [
				{
					"type": "location",
					"key": "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/locations/koreasouth"
				}
			]
		},
		"cloud_defender": {
			"enabled": true,
			"location_id": "defender-us-denver"
		},
		"credentials": [
			{
				"id": "E09F0AF8-18F8-49CE-B9AC-01C3E214B4EB",
				"purpose": "discover",
				"version": "2018-01-01"
			}
		],
		"status": {
			"status": "ok",
			"updated": 1493671342
		},
		"created": {
			"at": 1496673172,
			"by": "BBCAE827-22A6-433D-884E-22AABF2DC82B"
		},
		"modified": {
			"at": 1496673245,
			"by": "1443C74F-2AF7-4D7F-BF4A-FDFFE9C67182"
		},
		"version": 4
	}]`

	mux.HandleFunc(listDeploymentsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := []Deployment{
		{
			Version: 2,
			Status: Status{
				Status:  "ok",
				Updated: 1493671342,
			},
			Scope: Scope{
				Include: []Include{
					{
						Type:   "region",
						Key:    "/aws/us-east-1",
						Policy: Policy{ID: ""},
					},
					{
						Type:   "region",
						Key:    "/aws/us-east-2",
						Policy: Policy{ID: "A8E8B104-8F45-411D-A240-A30EA5FE25B0"},
					},
					{
						Type:   "vpc",
						Key:    "/aws/us-west-1/vpc/vpc-12345678",
						Policy: Policy{ID: "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"},
					},
				},
				Exclude: []Exclude{
					{
						Type: "region",
						Key:  "/aws/ap-southeast-1",
					},
				},
			},
			Scan: true,
			Platform: Platform{
				Type: "aws",
				ID:   "111111111111",
				Monitor: Monitor{
					Enabled:         true,
					CTInstallRegion: "us-east-1",
				},
				Default: false,
			},
			Name: "AWS Production Deployment",
			Modified: ModifiedCreated{
				At: 1496673245,
				By: "1443C74F-2AF7-4D7F-BF4A-FDFFE9C67182",
			},
			Mode:     "automatic",
			ID:       "AF50CF2D-3E77-46E3-B003-0262D77B2A65",
			Enabled:  true,
			Discover: true,
			Credentials: []Credential{
				{
					ID:      "E09F0AF8-18F8-49CE-B9AC-01C3E214B4EB",
					Purpose: "discover",
					Version: "2018-01-01",
				},
				{
					ID:      "D9B5066E-6145-4D41-9B3A-C123BDE64ABD",
					Purpose: "x-account-monitor",
					Version: "2018-01-01",
				},
			},
			Created: ModifiedCreated{
				At: 1496673172,
				By: "BBCAE827-22A6-433D-884E-22AABF2DC82B",
			},
			CloudDefender: CloudDefender{
				Enabled:    true,
				LocationID: "defender-us-denver",
			},
			AccountID: "01000001",
		},
		{
			Version: 4,
			Status: Status{
				Status:  "ok",
				Updated: 1493671342,
			},
			Scope: Scope{
				Include: []Include{
					{
						Type:   "location",
						Key:    "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/locations/centralus",
						Policy: Policy{ID: ""},
					},
					{
						Type:   "location",
						Key:    "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/locations/canadacentral",
						Policy: Policy{ID: "A8E8B104-8F45-411D-A240-A30EA5FE25B0"},
					},
					{
						Type:   "vnet",
						Key:    "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/resourceGroups/test-rg-01/providers/Microsoft.Network/virtualNetworks/test-vnet-01",
						Policy: Policy{ID: "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"},
					},
				},
				Exclude: []Exclude{
					{
						Type: "location",
						Key:  "/subscriptions/93B0107D-21AD-48B0-9101-C60F20E09966/locations/koreasouth",
					},
				},
			},
			Scan: true,
			Platform: Platform{
				Type: "azure",
				ID:   "93B0107D-21AD-48B0-9101-C60F20E09966",
				Monitor: Monitor{Enabled: false,
					CTInstallRegion: "",
				},
				Default: false,
			},
			Name: "Azure Production Deployment",
			Modified: ModifiedCreated{
				At: 1496673245,
				By: "1443C74F-2AF7-4D7F-BF4A-FDFFE9C67182",
			},
			Mode:     "manual",
			ID:       "86975343-2DA7-4E77-9F52-8488C4217191",
			Enabled:  true,
			Discover: true,
			Credentials: []Credential{
				{
					ID:      "E09F0AF8-18F8-49CE-B9AC-01C3E214B4EB",
					Purpose: "discover",
					Version: "2018-01-01",
				},
			},
			Created: ModifiedCreated{
				At: 1496673172,
				By: "BBCAE827-22A6-433D-884E-22AABF2DC82B",
			},
			CloudDefender: CloudDefender{
				Enabled:    true,
				LocationID: "defender-us-denver",
			},
			AccountID: "01000001",
		},
	}

	deployments, err := client.ListDeployments()
	if assert.NoError(t, err) {
		assert.Equal(t, deployments, want)
	}
}

func TestDeployments_ListDeploymentsMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listDeploymentsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.ListDeployments()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestDeployments_ListDeploymentsUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(listDeploymentsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.ListDeployments()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}

func TestDeployments_GetDeployment(t *testing.T) {
	setup()
	defer teardown()

	const response = `
	{
		"id": "AF50CF2D-3E77-46E3-B003-0262D77B2A65",
		"account_id": "01000001",
		"name": "AWS Production Deployment",
		"platform": {
			"type": "aws",
			"id": "111111111111",
			"monitor": {
				"enabled": true,
				"ct_install_region": "us-east-1"
			}
		},
		"mode": "automatic",
		"enabled": true,
		"discover": true,
		"scan": true,
		"scope": {
			"include": [
				{
					"type": "region",
					"key": "/aws/us-east-1"
				},
				{
					"type": "region",
					"key": "/aws/us-east-2",
					"policy": {
						"id": "A8E8B104-8F45-411D-A240-A30EA5FE25B0"
					}
				},
				{
					"type": "vpc",
					"key": "/aws/us-west-1/vpc/vpc-12345678",
					"policy": {
						"id": "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"
					}	
				}
			],
			"exclude": [
				{
					"type": "region",
					"key": "/aws/ap-southeast-1"
				}
			]
		},
		"cloud_defender": {
			"enabled": true,
			"location_id": "defender-us-denver"
		},
		"credentials": [
			{
				"id": "E09F0AF8-18F8-49CE-B9AC-01C3E214B4EB",
				"purpose": "discover",
				"version": "2018-01-01"
			},
			{
				"id": "D9B5066E-6145-4D41-9B3A-C123BDE64ABD",
				"purpose": "x-account-monitor",
				"version": "2018-01-01"
			}
		],
		"status": {
			"status": "ok",
			"updated": 1493671342
		},
		"created": {
			"at": 1496673172,
			"by": "BBCAE827-22A6-433D-884E-22AABF2DC82B"
		},
		"modified": {
			"at": 1496673245,
			"by": "1443C74F-2AF7-4D7F-BF4A-FDFFE9C67182"
		},
		"version": 2
	}`

	mux.HandleFunc(getDeploymentPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, response)
	})

	want := Deployment{
		Version: 2,
		Status: Status{
			Status:  "ok",
			Updated: 1493671342,
		},
		Scope: Scope{
			Include: []Include{
				{
					Type:   "region",
					Key:    "/aws/us-east-1",
					Policy: Policy{ID: ""},
				},
				{
					Type:   "region",
					Key:    "/aws/us-east-2",
					Policy: Policy{ID: "A8E8B104-8F45-411D-A240-A30EA5FE25B0"},
				},
				{
					Type:   "vpc",
					Key:    "/aws/us-west-1/vpc/vpc-12345678",
					Policy: Policy{ID: "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"},
				},
			},
			Exclude: []Exclude{
				{
					Type: "region",
					Key:  "/aws/ap-southeast-1",
				},
			},
		},
		Scan: true,
		Platform: Platform{
			Type: "aws",
			ID:   "111111111111",
			Monitor: Monitor{
				Enabled:         true,
				CTInstallRegion: "us-east-1",
			},
			Default: false,
		},
		Name: "AWS Production Deployment",
		Modified: ModifiedCreated{
			At: 1496673245,
			By: "1443C74F-2AF7-4D7F-BF4A-FDFFE9C67182",
		},
		Mode:     "automatic",
		ID:       "AF50CF2D-3E77-46E3-B003-0262D77B2A65",
		Enabled:  true,
		Discover: true,
		Credentials: []Credential{
			{
				ID:      "E09F0AF8-18F8-49CE-B9AC-01C3E214B4EB",
				Purpose: "discover",
				Version: "2018-01-01",
			},
			{
				ID:      "D9B5066E-6145-4D41-9B3A-C123BDE64ABD",
				Purpose: "x-account-monitor",
				Version: "2018-01-01",
			},
		},
		Created: ModifiedCreated{
			At: 1496673172,
			By: "BBCAE827-22A6-433D-884E-22AABF2DC82B",
		},
		CloudDefender: CloudDefender{
			Enabled:    true,
			LocationID: "defender-us-denver",
		},
		AccountID: "01000001",
	}

	deployment, err := client.GetDeployment(testDeploymentId)
	if assert.NoError(t, err) {
		assert.Equal(t, deployment, want)
	}
}

func TestDeployments_GetDeploymentMakeRequestError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getDeploymentPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetDeployment(testDeploymentId)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error from makeRequest: HTTP status 401: invalid credentials")
}

func TestDeployments_GetDeploymentUnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(getDeploymentPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, "not json")
	})

	_, err := client.GetDeployment(testDeploymentId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), testUnmarshalError)
}
