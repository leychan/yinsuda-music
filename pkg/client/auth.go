package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TokenProvider handles fetching and refreshing access tokens.
type TokenProvider struct {
	appId      string
	appSecret  string
	authUrl    string // e.g., "https://api.yinsuda.com/oauth2/token"
	httpClient *http.Client

	lock        sync.RWMutex
	accessToken string
	expiresAt   time.Time
}

func NewTokenProvider(appId, appSecret, authUrl string, client *http.Client) *TokenProvider {
	if client == nil {
		client = http.DefaultClient
	}
	return &TokenProvider{
		appId:      appId,
		appSecret:  appSecret,
		authUrl:    authUrl,
		httpClient: client,
	}
}

// GetAccessToken returns a valid access token, refreshing if necessary.
func (p *TokenProvider) GetAccessToken() (string, error) {
	p.lock.RLock()
	token := p.accessToken
	expiry := p.expiresAt
	p.lock.RUnlock()

	// Check if token is valid (with 30s buffer)
	if token != "" && time.Now().Add(30*time.Second).Before(expiry) {
		return token, nil
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	// Double check after acquiring write lock
	if p.accessToken != "" && time.Now().Add(30*time.Second).Before(p.expiresAt) {
		return p.accessToken, nil
	}

	newToken, expireSeconds, err := p.fetchToken()
	if err != nil {
		return "", err
	}

	p.accessToken = newToken
	p.expiresAt = time.Now().Add(time.Duration(expireSeconds) * time.Second)
	return newToken, nil
}

func (p *TokenProvider) fetchToken() (string, int, error) {
	reqBody := map[string]string{
		"appId":     p.appId,
		"appSecret": p.appSecret,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequest("POST", p.authUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("auth request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", 0, err
	}

	if !tokenResp.Success || tokenResp.Code != "0" {
		return "", 0, fmt.Errorf("auth failed: %s (%s)", tokenResp.Message, tokenResp.Msg)
	}

	return tokenResp.Data.AccessToken, tokenResp.Data.Expire, nil
}
