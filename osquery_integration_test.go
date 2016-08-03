package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestIntegrationEnrollHostBadSecret(t *testing.T) {
	var req IntegrationRequests
	req.New(t)

	config.Osquery.EnrollSecret = "super secret"

	// Check that a bad enroll secret causes the appropriate error code and
	// error JSON to be returned

	resp := req.EnrollHost("bad secret", "fake_uuid")

	if resp.Code != http.StatusUnauthorized {
		t.Error("Should error with invalid enroll secret")
	}

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		req.t.Fatalf("JSON decode error: %s JSON contents:\n %s", err.Error(), resp.Body.Bytes())
	}

	if body["error"] != "Invalid enroll secret" {
		req.t.Errorf("Incorrect/missing error message: %s", body["error"])
	}

	if invalid, ok := body["node_invalid"]; ok && invalid != true {
		req.t.Errorf("Expected node_invalid = true")
	}
}

func TestIntegrationEnrollHostMissingIdentifier(t *testing.T) {
	var req IntegrationRequests
	req.New(t)

	config.Osquery.EnrollSecret = "super secret"

	// Check that an empty host identifier causes the appropriate error code and
	// error JSON to be returned

	resp := req.EnrollHost("super secret", "")

	if resp.Code != http.StatusBadRequest {
		t.Error("Should error with missing host identifier")
	}

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		req.t.Fatalf("JSON decode error: %s JSON contents:\n %s", err.Error(), resp.Body.Bytes())
	}

	if body["error"] != "Missing host identifier" {
		req.t.Errorf("Incorrect/missing error message: %s", body["error"])
	}

	if invalid, ok := body["node_invalid"]; ok && invalid != true {
		req.t.Errorf("Expected node_invalid = true")
	}
}

func TestIntegrationEnrollHostGood(t *testing.T) {
	var req IntegrationRequests
	req.New(t)

	config.Osquery.EnrollSecret = "super secret"

	// Make a good request and verify that a node key is returned. Also
	// check that the DB recorded the information.

	resp := req.EnrollHost("super secret", "fake_host_1")

	if resp.Code != http.StatusOK {
		t.Error("Status should be StatusOK")
	}

	t.Logf("Response body:\n%s", string(resp.Body.Bytes()))

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		req.t.Fatalf("JSON decode error: %s JSON contents:\n %s", err.Error(), resp.Body.Bytes())
	}

	if _, ok := body["error"]; ok {
		req.t.Errorf("Unexpected error message: %s", body["error"])
	}

	if invalid, ok := body["node_invalid"]; ok && invalid == true {
		req.t.Errorf("Expected node_invalid = false")
	}

	nodeKey, ok := body["node_key"]
	if !ok || nodeKey == "" {
		req.t.Errorf("Expected node_key")
	}

	var host Host
	err = req.db.Where("uuid = ?", "fake_host_1").First(&host).Error
	if err != nil {
		req.t.Fatal("Host not saved to DB: %s", err.Error())
	}

	if host.NodeKey != nodeKey {
		req.t.Error("Saved node key different than response key, %s != %s",
			host.NodeKey, nodeKey)
	}

	// Enroll again and check that node key changes

	resp = req.EnrollHost("super secret", "fake_host_1")

	if resp.Code != http.StatusOK {
		t.Error("Status should be StatusOK")
	}

	t.Logf("Response body:\n%s", string(resp.Body.Bytes()))

	body = map[string]interface{}{}
	err = json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		req.t.Fatalf("JSON decode error: %s JSON contents:\n %s", err.Error(), resp.Body.Bytes())
	}

	if _, ok := body["error"]; ok {
		req.t.Errorf("Unexpected error message: %s", body["error"])
	}

	if invalid, ok := body["node_invalid"]; ok && invalid == true {
		req.t.Errorf("Expected node_invalid = false")
	}

	newNodeKey, ok := body["node_key"]
	if !ok || nodeKey == "" {
		req.t.Errorf("Expected node_key")
	}

	if newNodeKey == nodeKey {
		req.t.Errorf("Node key should have changed, %s == %s", newNodeKey, nodeKey)
	}

}
