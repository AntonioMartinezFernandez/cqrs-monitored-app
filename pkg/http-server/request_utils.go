package http_server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const reverseProxyForwardedByHeader = "X-Forwarded-For"

func ClientIp(req *http.Request) string {
	ipAddress := req.RemoteAddr
	fwdAddress := req.Header.Get(reverseProxyForwardedByHeader)
	if fwdAddress != "" {
		ipAddress = fwdAddress

		ips := strings.Split(fwdAddress, ", ")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}

	return ipAddress
}

func CloneRequest(r *http.Request) *http.Request {
	var bodyBytes []byte
	newRequest := *r.WithContext(r.Context())

	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	newRequest.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return &newRequest
}

func AllParamsRequest(r *http.Request) (map[string]any, error) {
	newRequest := CloneRequest(r)

	allParams, err := ConvertRequestToBodyMap(newRequest)
	if err != nil {
		return allParams, err
	}

	queryParams := QueryParams(newRequest)

	for k, v := range queryParams {
		allParams[k] = v
	}

	return allParams, nil
}

func ConvertRequestToBodyMap(r *http.Request) (map[string]any, error) {
	requestBody := make(map[string]any)
	var err error
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return requestBody, err
	}
	if err = json.Unmarshal(b, &requestBody); err != nil {
		return requestBody, err
	}

	return requestBody, err
}

func QueryParams(r *http.Request) map[string]any {
	query := make(map[string]any)

	for k, v := range r.URL.Query() {
		query[k] = v[0]
	}

	return query
}
