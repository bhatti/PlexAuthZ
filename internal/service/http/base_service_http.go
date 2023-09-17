package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"io"
	"time"
)

// baseHTTPClient - base http client
type baseHTTPClient struct {
	client  web.HTTPClient
	baseURL string
}

// post helper
func (h *baseHTTPClient) post(
	ctx context.Context,
	path string,
	req any,
	res any,
) (status int, respHeaders map[string]string, err error) {
	return h.handle(ctx, path, "POST", make(map[string]string), req, res)
}

// put helper
func (h *baseHTTPClient) put(
	ctx context.Context,
	path string,
	req any,
	res any,
) (status int, respHeaders map[string]string, err error) {
	return h.handle(ctx, path, "PUT", make(map[string]string), req, res)
}

// get helper
func (h *baseHTTPClient) get(
	ctx context.Context,
	path string,
	params map[string]string,
	res any,
) (status int, respHeaders map[string]string, err error) {
	return h.handle(ctx, path, "GET", params, nil, res)
}

// del helper
func (h *baseHTTPClient) del(
	ctx context.Context,
	path string,
) (status int, respHeaders map[string]string, err error) {
	return h.handle(ctx, path, "DELETE", make(map[string]string), nil, nil)
}

// handle helper
func (h *baseHTTPClient) handle(
	ctx context.Context,
	path string,
	method string,
	params map[string]string,
	req any,
	res any,
) (status int, respHeaders map[string]string, err error) {
	var reqBody io.ReadCloser
	var resBody io.ReadCloser
	started := time.Now()
	if req != nil {
		b, err := json.Marshal(req)
		if err != nil {
			return 0, nil, err
		}
		reqBody = io.NopCloser(bytes.NewReader(b))
	}
	status, _, resBody, respHeaders, err = h.client.Handle(
		ctx,
		h.baseURL+path,
		method,
		map[string]string{"Content-Type": "application/json"},
		params,
		reqBody,
	)
	elapsed := time.Since(started).String()
	if err != nil {
		return status, nil, err
	}
	if status >= 300 {
		return status, respHeaders,
			fmt.Errorf("failed to invoke %s %s with %v due to status %d, time: %s [%s]",
				method, path, req, status, elapsed, domain.NetworkCode)
	}
	if resBody != nil && res != nil {
		b, err := io.ReadAll(resBody)
		if err != nil {
			return status, respHeaders, err
		}
		return status, respHeaders, json.Unmarshal(b, res)
	}
	return status, respHeaders, nil
}
