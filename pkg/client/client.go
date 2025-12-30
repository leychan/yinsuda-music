package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Client is the main entry point for interacting with the Yinsuda Music API.
type Client struct {
	appId        string
	tokenProvider *TokenProvider
	httpClient   *http.Client
	baseUrl      string
}

// NewClient creates a new Yinsuda Music API client.
// baseUrl should be the root URL of the API, e.g., "https://api.yinsuda.com"
func NewClient(appId, appSecret, baseUrl string) *Client {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	// Assume the auth endpoint is relative to baseUrl, e.g. /oauth2/token
	authUrl := fmt.Sprintf("%s/oauth2/token", baseUrl)
	
	return &Client{
		appId:        appId,
		tokenProvider: NewTokenProvider(appId, appSecret, authUrl, httpClient),
		httpClient:   httpClient,
		baseUrl:      baseUrl,
	}
}

// Do performs a request to the API, handling authentication and signing.
// method: GET, POST, etc.
// path: relative path, e.g., "/musician/list"
// body: request body struct (will be marshaled to JSON), or nil
// result: pointer to struct where response data will be unmarshaled
func (c *Client) Do(method, path string, query url.Values, body interface{}, result interface{}) error {
	// 1. Get Access Token
	accessToken, err := c.tokenProvider.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// 2. Prepare Body
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %w", err)
		}
	}

	// 3. Prepare Sign Params
	timestamp := time.Now().Format("20060102150405")
	traceId := fmt.Sprintf("musician-openapi_%s", uuid.New().String())
	
	signParams := SignParams{
		AppId:       c.appId,
		AccessToken: accessToken,
		Timestamp:   timestamp,
		SignMethod:  "md5",
		TraceId:     traceId,
		// SignVersion excluded from calculation based on Java ref
	}

	// 4. Calculate Sign
	// Note: We need to pass the query values separately if it's a GET request or has query params
	// path should strictly be the path, e.g., /foo/bar.
	sign := CalculateSign(signParams, bodyBytes, path, query, c.tokenProvider.appSecret)

	// 5. Construct Request
	fullUrl := fmt.Sprintf("%s%s", c.baseUrl, path)
	if len(query) > 0 {
		fullUrl += "?" + query.Encode()
	}

	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 6. Set Headers
	q := req.Header
	q.Set("Content-Type", "application/json")
	q.Set("appId", c.appId)
	q.Set("accessToken", accessToken)
	q.Set("timestamp", timestamp)
	q.Set("signMethod", "md5")
	q.Set("traceId", traceId)
	q.Set("sign", sign)
	q.Set("signVersion", "v2")
	// "source" is optional, not setting it for now.

	// 7. Execute
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyDump, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyDump))
	}

	// 8. Parse Response
	// If result is expected
	if result != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		// It seems the API always returns a standard envelope.
		// We can try to decode strictly if the user provided a struct matching the inner data,
		// or matching the full envelope. Data models usually just want the 'data' part,
		// but checking 'code' is important.
		
		// Let's decode into BaseResponse first to check generic errors.
		var baseResp BaseResponse
		if err := json.Unmarshal(respBody, &baseResp); err != nil {
			return fmt.Errorf("failed to parse base response: %w", err)
		}

		if !baseResp.Success || baseResp.Code != 0 {
			return fmt.Errorf("API error: %s (%s) %s", baseResp.Message, baseResp.Msg, string(respBody))
		}

		// Now unmarshal into the specific result
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal into result: %w", err)
		}
	}

	return nil
}
