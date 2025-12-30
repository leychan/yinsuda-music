package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Integration(t *testing.T) {
	// 1. Setup Mock Server
	mux := http.NewServeMux()
	
	// Mock Auth Endpoint
	mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Auth expected POST, got %s", r.Method)
		}
		
		// Verify payload
		var req map[string]string
		json.NewDecoder(r.Body).Decode(&req)
		if req["appId"] != "testAppId" || req["appSecret"] != "testAppSecret" {
			t.Errorf("Auth params invalid: %v", req)
		}

		resp := TokenResponse{
			Code: "0",
			Success: true,
			Data: TokenData{
				AccessToken: "mock_access_token",
				Expire:      3600,
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Mock Business Endpoint
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// Verify common headers
		if r.Header.Get("appId") != "testAppId" {
			t.Errorf("Missing/Invalid appId header")
		}
		if r.Header.Get("accessToken") != "mock_access_token" {
			t.Errorf("Missing/Invalid accessToken header")
		}
		if r.Header.Get("sign") == "" {
			t.Errorf("Missing sign header")
		}
		
		// Return dummy response
		resp := BaseResponse{
			Code:    0,
			Success: true,
			Data:    map[string]string{"foo": "bar"},
		}
		json.NewEncoder(w).Encode(resp)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// 2. Init Client
	// Pass the server URL as baseUrl.
	// Note: Our client appends /oauth2/token, so server URL is correct base.
	client := NewClient("testAppId", "testAppSecret", server.URL)

	// 3. Call Business API
	var result BaseResponse
	err := client.Do("POST", "/api/data", nil, map[string]string{"input": "val"}, &result)
	if err != nil {
		t.Fatalf("Client.Do failed: %v", err)
	}

	// 4. Verify Result
	dataMap, ok := result.Data.(map[string]interface{})
	if !ok || dataMap["foo"] != "bar" {
		t.Errorf("Unexpected result data: %v", result.Data)
	}
}
