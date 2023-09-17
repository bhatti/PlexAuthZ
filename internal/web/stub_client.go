package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// StubHTTPResponse defines stub response
type StubHTTPResponse struct {
	Filename      string
	Bytes         []byte
	Status        int
	Headers       map[string][]string
	Error         error
	sleepDuration time.Duration
}

// NewStubHTTPResponseError creates stubbed response with error
func NewStubHTTPResponseError(status int, sleep time.Duration, err error) *StubHTTPResponse {
	return &StubHTTPResponse{Status: status, sleepDuration: sleep, Error: err,
		Headers: map[string][]string{"Content-Type": {"application/json"}}}
}

// NewStubHTTPResponse creates stubbed response
func NewStubHTTPResponse(status int, unk any) *StubHTTPResponse {
	res := &StubHTTPResponse{Status: status,
		Headers: map[string][]string{"Content-Type": {"application/json"}}}
	if unk == nil {
		return res
	}
	switch unk.(type) {
	case string:
		if _, err := os.Stat(unk.(string)); err == nil {
			res.Filename = unk.(string)
		} else {
			res.Bytes = []byte(unk.(string))
		}
	default:
		b, err := json.Marshal(unk)
		if err != nil {
			panic(fmt.Errorf("failed to serialize %v due to error %w", unk, err))
		}
		if _, err := os.Stat(string(b)); err == nil {
			res.Filename = string(b)
		} else {
			res.Bytes = b
		}
	}
	return res
}

// WithHeader sets header.
func (r *StubHTTPResponse) WithHeader(name string, val string) *StubHTTPResponse {
	r.Headers[name] = []string{val}
	return r
}

// StubHTTPClient implements HTTPClient for stubbed response
type StubHTTPClient struct {
	mappingByMethodURL map[string]*StubHTTPResponse
}

// NewStubHTTPClient - creates structure for HTTPClient
func NewStubHTTPClient() *StubHTTPClient {
	return &StubHTTPClient{
		mappingByMethodURL: make(map[string]*StubHTTPResponse),
	}
}

// AddMapping adds mapping for stub response
func (w *StubHTTPClient) AddMapping(method string, url string, resp *StubHTTPResponse) {
	w.mappingByMethodURL[method+url] = resp
}

// getMapping finds mapping for stub response
func (w *StubHTTPClient) getMapping(method string, url string) *StubHTTPResponse {
	return w.mappingByMethodURL[method+url]
}

// Handle makes HTTP request
func (w *StubHTTPClient) Handle(
	_ context.Context,
	url string,
	method string,
	_ map[string][]string,
	_ map[string]string,
	_ io.ReadCloser) (int, string, io.ReadCloser, map[string][]string, error) {
	if url == "" {
		return 500, "", nil, nil, fmt.Errorf("url is not specified")
	}
	if method == "" {
		return 500, "", nil, nil, fmt.Errorf("method is not specified")
	}
	log.WithFields(log.Fields{"component": "stub-web", "url": url, "method": method}).Debugf("BEGIN")
	resp := w.getMapping(method, url)
	if resp == nil {
		return 404, "", nil, nil, fmt.Errorf("couldn't find URL '%s' method '%s' in mapping: %v",
			url, method, w.mappingByMethodURL)
	}
	if resp.sleepDuration > 0 {
		time.Sleep(resp.sleepDuration)
	}
	if len(resp.Bytes) > 0 {
		return resp.Status, "", io.NopCloser(bytes.NewReader(resp.Bytes)), resp.Headers, resp.Error
	}
	if resp.Error != nil {
		return resp.Status, "", nil, resp.Headers, resp.Error
	}
	b, err := os.ReadFile(resp.Filename)
	if err != nil {
		return 404, "", nil, resp.Headers, fmt.Errorf("error reading file %v for url %v due to %w", resp.Filename, url, err)
	}
	return resp.Status, "", io.NopCloser(bytes.NewReader(b)), resp.Headers, resp.Error
}
