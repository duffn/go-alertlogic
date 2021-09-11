package alertlogic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	modifyExternalDNSNameAssetPath = fmt.Sprintf("/%s/%s/deployments/%s/assets", assetsWriteServicePath, testAccountId, testDeploymentId)
)

func TestAssetsWrite_CreateExternalDNSNameAsset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(modifyExternalDNSNameAssetPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
	})

	assetResponse, err := client.CreateExternalDNSNameAsset(testDeploymentId, "abcd-1234.elb.us-east-1.amazonaws.com")

	if assert.NoError(t, err) {
		assert.Equal(t, assetResponse, http.StatusCreated)
	}
}

func TestAssetsWrite_CreateExternalDNSNameAssetError(t *testing.T) {
	setup()
	defer teardown()

	errorResponse := `{"error": "Invalid operation"}`

	mux.HandleFunc(modifyExternalDNSNameAssetPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, errorResponse)
	})

	respCode, err := client.CreateExternalDNSNameAsset(testDeploymentId, "qwert-9876.elb.us-east-1.amazonaws.com")

	assert.Error(t, err)
	assert.Equal(t, respCode, http.StatusBadRequest)
	assert.Equal(t, err.Error(), fmt.Sprintf("error from makeRequest: %s", errorResponse))
}

func TestAssetsWrite_UpdateExternalDNSNameAsset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(modifyExternalDNSNameAssetPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
	})

	assetResponse, err := client.UpdateExternalDNSNameAsset(testDeploymentId, "qwert-9876.elb.us-east-1.amazonaws.com", "abcd-1234.elb.us-east-1.amazonaws.com")

	if assert.NoError(t, err) {
		assert.Equal(t, assetResponse, http.StatusCreated)
	}
}

func TestAssetsWrite_UpdateExternalDNSNameAssetError(t *testing.T) {
	setup()
	defer teardown()

	errorResponse := `{"error": "Invalid operation"}`

	mux.HandleFunc(modifyExternalDNSNameAssetPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, errorResponse)
	})

	respCode, err := client.UpdateExternalDNSNameAsset(testDeploymentId, "qwert-9876.elb.us-east-1.amazonaws.com", "abcd-1234.elb.us-east-1.amazonaws.com")

	assert.Error(t, err)
	assert.Equal(t, respCode, http.StatusBadRequest)
	assert.Equal(t, err.Error(), fmt.Sprintf("error from makeRequest: %s", errorResponse))
}

func TestAssetsWrite_RemoveExternalDNSNameAsset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(modifyExternalDNSNameAssetPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	assetResponse, err := client.RemoveExternalDNSNameAsset(testDeploymentId, "qwert-9876.elb.us-east-1.amazonaws.com")

	if assert.NoError(t, err) {
		assert.Equal(t, assetResponse, http.StatusNoContent)
	}
}

func TestAssetsWrite_RemoveExternalDNSNameAssetError(t *testing.T) {
	setup()
	defer teardown()

	errorResponse := `{"error": "Invalid operation"}`

	mux.HandleFunc(modifyExternalDNSNameAssetPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, errorResponse)
	})

	respCode, err := client.RemoveExternalDNSNameAsset(testDeploymentId, "qwert-9876.elb.us-east-1.amazonaws.com")

	assert.Error(t, err)
	assert.Equal(t, respCode, http.StatusBadRequest)
	assert.Equal(t, err.Error(), fmt.Sprintf("error from makeRequest: %s", errorResponse))
}
