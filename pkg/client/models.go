package client

// BaseResponse represents the standard JSON response structure.
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"traceId"`
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
}

// TokenData represents the data field in the token response.
type TokenData struct {
	AccessToken string `json:"accessToken"`
	Expire      int    `json:"expire"` // Expiration time in seconds? Or timestamp? Prompt says "expire": 900, likely seconds.
}

// TokenResponse represents the full response for /oauth2/token.
type TokenResponse struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    TokenData `json:"data"`
	TraceId string    `json:"traceId"`
	Success bool      `json:"success"`
	Msg     string    `json:"msg"`
}
