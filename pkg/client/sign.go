package client

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// SignParams holds the public parameters required for signature calculation.
// Note: "sign" itself is excluded.
type SignParams struct {
	AppId       string
	AccessToken string
	Timestamp   string
	SignMethod  string
	TraceId     string
	// SignVersion field is excluded from signature in Java implementation
	Source string
}

// CalculateSign generates the MD5 signature based on the API specification.
// Algorithm:
// (a) Filter headers (appId, accessToken, timestamp, signMethod, traceId, source).
//
//	Sort by key asc, concat as key:value; (skip empty values).
//
// (b) md5(body)
// (c) md5(urlPath + query)
// (d) md5(appSecret)
// (e) md5(headerMd5 + bodyMd5 + urlMd5 + appSecretMd5) -> hex
func CalculateSign(params SignParams, body []byte, urlPath string, query url.Values, appSecret string) string {
	// (a) Public params
	// Java implementation: appId, accessToken, timestamp, signMethod, traceId, source
	m := map[string]string{
		"appId":       params.AppId,
		"accessToken": params.AccessToken,
		"timestamp":   params.Timestamp,
		"signMethod":  params.SignMethod,
		"traceId":     params.TraceId,
		"source":      params.Source,
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		// Java: if (entry.getValue() != null && !"".equals(entry.getValue()))
		if m[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var headerBuilder strings.Builder
	for _, k := range keys {
		headerBuilder.WriteString(fmt.Sprintf("%s:%s;", k, m[k]))
	}
	headerStr := headerBuilder.String()
	headerMd5 := md5Hash([]byte(headerStr))

	// (b) Body
	// If body is empty or nil, we still hash it (empty string hash)
	bodyMd5 := md5Hash(body)

	// (c) URL
	// "包含query参数(如有)，不包含域名，比如/foo/bar?key=value"
	// We need to ensure the query is formatted correctly if present.
	// Standard url.Values.Encode() sorts keys.
	fullPath := urlPath
	if len(query) > 0 {
		fullPath = fmt.Sprintf("%s?%s", urlPath, query.Encode())
	}
	urlMd5 := md5Hash([]byte(fullPath))

	// (d) AppSecret
	appSecretMd5 := md5Hash([]byte(appSecret))

	// (e) Final Sign
	// buffer: headerMd5 (bytes) + bodyMd5 (bytes) ...
	// The prompt says "拼接md5字节数组" (concatenate md5 byte arrays).
	// md5Hash returns []byte.
	finalInput := make([]byte, 0, 16*4)
	finalInput = append(finalInput, headerMd5...)
	finalInput = append(finalInput, bodyMd5...)
	finalInput = append(finalInput, urlMd5...)
	finalInput = append(finalInput, appSecretMd5...)

	finalMd5 := md5.Sum(finalInput)
	return hex.EncodeToString(finalMd5[:])
}

// md5Hash returns the raw 16 bytes MD5 hash
func md5Hash(data []byte) []byte {
	h := md5.Sum(data)
	return h[:]
}
