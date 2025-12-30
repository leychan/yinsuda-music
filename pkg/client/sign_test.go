package client

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"
)

func TestCalculateSign(t *testing.T) {
	// Setup test data
	params := SignParams{
		AppId:       "123456",
		AccessToken: "token123",
		Timestamp:   "20210101120000",
		SignMethod:  "md5",
		TraceId:     "trace123",
		// SignVersion removed
	}
	appSecret := "secret123"
	body := []byte(`{"key":"value"}`)
	urlPath := "/test/api"
	// Empty query for this test
	query := url.Values{}

	// Calculate manually
	// (a) Header
	// Sorted keys: accessToken, appId, signMethod, timestamp, traceId
	// Note: signVersion is REMOVED in new logic.
	headerStr := "accessToken:token123;appId:123456;signMethod:md5;timestamp:20210101120000;traceId:trace123;"
	headerMd5 := md5Hash([]byte(headerStr))

	// (b) Body
	bodyMd5 := md5Hash(body)

	// (c) URL
	urlMd5 := md5Hash([]byte(urlPath))

	// (d) AppSecret
	secretMd5 := md5Hash([]byte(appSecret))

	// (e) Final
	finalInput := make([]byte, 0)
	finalInput = append(finalInput, headerMd5...)
	finalInput = append(finalInput, bodyMd5...)
	finalInput = append(finalInput, urlMd5...)
	finalInput = append(finalInput, secretMd5...)
	finalSum := md5.Sum(finalInput)
	expectedSign := hex.EncodeToString(finalSum[:])

	// Execute
	actualSign := CalculateSign(params, body, urlPath, query, appSecret)

	if actualSign != expectedSign {
		t.Errorf("CalculateSign mismatch.\nExpected: %s\nActual:   %s", expectedSign, actualSign)
		
		fmt.Printf("HeaderStr: %s\n", headerStr)
		fmt.Printf("HeaderMD5: %x\n", headerMd5)
		fmt.Printf("BodyMD5:   %x\n", bodyMd5)
		fmt.Printf("UrlMD5:    %x\n", urlMd5)
		fmt.Printf("SecretMD5: %x\n", secretMd5)
	}
}

func TestCalculateSign_WithQuery(t *testing.T) {
	// Setup test data
	params := SignParams{
		AppId:       "123456",
		AccessToken: "token123",
		Timestamp:   "20210101120000",
		SignMethod:  "md5",
		TraceId:     "trace123",
		// SignVersion removed
	}
	appSecret := "secret123"
	body := []byte("") // Empty body
	urlPath := "/test/api"
	query := url.Values{}
	query.Set("b", "2")
	query.Set("a", "1")

	// Verify query sorting logic inside CalculateSign
	// query.Encode() produces "a=1&b=2"
	
	// (a) Header (same as above)
	headerStr := "accessToken:token123;appId:123456;signMethod:md5;timestamp:20210101120000;traceId:trace123;"
	headerMd5 := md5Hash([]byte(headerStr))

	// (b) Body
	bodyMd5 := md5Hash(body)

	// (c) URL
	// /test/api?a=1&b=2
	fullPath := "/test/api?a=1&b=2"
	urlMd5 := md5Hash([]byte(fullPath))

	// (d) AppSecret
	secretMd5 := md5Hash([]byte(appSecret))

	// (e) Final
	finalInput := make([]byte, 0)
	finalInput = append(finalInput, headerMd5...)
	finalInput = append(finalInput, bodyMd5...)
	finalInput = append(finalInput, urlMd5...)
	finalInput = append(finalInput, secretMd5...)
	finalSum := md5.Sum(finalInput)
	expectedSign := hex.EncodeToString(finalSum[:])

	// Execute
	actualSign := CalculateSign(params, body, urlPath, query, appSecret)

	if actualSign != expectedSign {
		t.Errorf("CalculateSign (Query) mismatch.\nExpected: %s\nActual:   %s", expectedSign, actualSign)
	}
}
