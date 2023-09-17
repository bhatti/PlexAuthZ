package web

import (
	"bytes"
	"context"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

const todoURL = "https://jsonplaceholder.typicode.com/todos"

func Test_ShouldNotRealGetWithBadUrl(t *testing.T) {
	w := newTestNewHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		"uuu",
		"mmm",
		map[string]string{"key": "value"},
		make(map[string]string),
		nil)
	if err != nil {
		t.Logf("unexpected response GetByIDs error " + err.Error())
	}
}

func Test_ShouldNotRealGetWithBadMethod(t *testing.T) {
	w := newTestNewHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL+"/1",
		"mmm",
		map[string]string{"key": "value"},
		make(map[string]string),
		nil)
	if err != nil {
		t.Logf("unexpected response GetByIDs error " + err.Error())
	}
}

func Test_ShouldRealGet(t *testing.T) {
	w := newTestNewHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL+"/1",
		"GET",
		map[string]string{"key": "value"},
		make(map[string]string),
		nil)
	if err != nil {
		t.Logf("unexpected response GetByIDs error " + err.Error())
	}
}

func Test_ShouldRealDelete(t *testing.T) {
	w := newTestNewHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL+"/1",
		"DELETE",
		map[string]string{"key": "value"},
		make(map[string]string),
		nil)
	if err != nil {
		t.Logf("unexpected response Delete error " + err.Error())
	}
}

func Test_ShouldRealDeleteBody(t *testing.T) {
	w := newTestNewHTTPClient()
	body := io.NopCloser(bytes.NewReader([]byte("hello")))
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL+"/1",
		"DELETE",
		map[string]string{"key": "value"},
		make(map[string]string),
		body)
	if err != nil {
		t.Logf("unexpected response Delete error " + err.Error())
	}
}

func Test_ShouldRealPostError(t *testing.T) {
	w := newTestNewHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL+"_____",
		"POST",
		map[string]string{"key": "value"},
		map[string]string{},
		nil)
	if err == nil {
		t.Logf("expected response Post error ")
	}
}

func Test_ShouldRealPost(t *testing.T) {
	w := newTestNewHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL,
		"POST",
		map[string]string{"key": "value"},
		map[string]string{},
		nil)
	if err != nil {
		t.Logf("unexpected response Post error " + err.Error())
	}
}

func Test_ShouldRealPostForm(t *testing.T) {
	w := newTestNewHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL,
		"POST",
		map[string]string{"name": "value"},
		map[string]string{"name": "value"},
		nil)
	if err != nil {
		t.Logf("unexpected response Post error " + err.Error())
	}
}

func Test_ShouldRealPostBody(t *testing.T) {
	w := newTestNewHTTPClient()
	body := io.NopCloser(bytes.NewReader([]byte("hello")))
	_, _, _, _, err := w.Handle(
		context.Background(),
		todoURL,
		"POST",
		map[string]string{"name": "value"},
		map[string]string{},
		body)
	if err != nil {
		t.Logf("unexpected response Post error " + err.Error())
	}
}

func Test_ShouldNotGetRemoteIPAddressFromURL(t *testing.T) {
	require.Equal(t, "", getRemoteIPAddressFromURL("xxx"))
	require.Contains(t, getRemoteIPAddressFromURL("http://localhost"), "127.0.0.1")
	require.True(t, len(getLocalIPAddresses()) > 0)
}

func Test_ShouldGetgetProxyEnv(t *testing.T) {
	require.Equal(t, 3, len(getProxyEnv()))
}

func Test_ShouldNotExecuteHttpClientWithNilRequest(t *testing.T) {
	config := &domain.Config{ProxyURL: "http://localhost:8000"}
	status, _, _, _, err := NewHTTPClient(config).execute(nil, nil, nil)
	require.Error(t, err)
	require.Equal(t, 500, status)
}

func Test_ShouldGetHttpClientWithProxy(t *testing.T) {
	config := &domain.Config{ProxyURL: "xyz"}
	require.NotNil(t, httpClient(config))
	config = &domain.Config{ProxyURL: "ftp://localhost:8000"}
	require.NotNil(t, httpClient(config))
	config = &domain.Config{ProxyURL: "http://localhost:8000"}
	require.NotNil(t, httpClient(config))
}

func newTestNewHTTPClient() *DefaultHTTPClient {
	c := domain.Config{}
	return NewHTTPClient(&c)
}
