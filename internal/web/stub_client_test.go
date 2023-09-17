package web

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func Test_ShouldNotStubGetWithBadUrl(t *testing.T) {
	w := NewStubHTTPClient()
	_, _, _, _, err := w.Handle(
		context.Background(),
		"uuu",
		"mmm",
		map[string][]string{"key": {"value"}},
		make(map[string]string),
		nil)
	require.Error(t, err)
}

func Test_ShouldNotPostStubClientWithoutURL(t *testing.T) {
	w := NewStubHTTPClient()
	w.AddMapping("POST", "/path1", NewStubHTTPResponse(200, "test body"))
	_, _, _, _, err := w.Handle(
		context.Background(),
		"",
		"POST",
		map[string][]string{"key": {"value"}},
		make(map[string]string),
		nil)
	require.Error(t, err)
}

func Test_ShouldNotPostStubClientWithoutMethod(t *testing.T) {
	w := NewStubHTTPClient()
	w.AddMapping("POST", "/path1", NewStubHTTPResponse(200, "test body"))
	_, _, _, _, err := w.Handle(
		context.Background(),
		"/path",
		"",
		map[string][]string{"key": {"value"}},
		make(map[string]string),
		nil)
	require.Error(t, err)
}

func Test_ShouldPostStubClientWithIntegerResponse(t *testing.T) {
	w := NewStubHTTPClient()
	w.AddMapping("POST", "/path1", NewStubHTTPResponse(200, 3))
	status, _, reader, _, err := w.Handle(
		context.Background(),
		"/path1",
		"POST",
		map[string][]string{"key": {"value"}},
		make(map[string]string),
		nil)
	require.NoError(t, err)
	require.Equal(t, 200, status)
	b, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, 1, len(b))
}

func Test_ShouldPostStubClientWithStringResponse(t *testing.T) {
	w := NewStubHTTPClient()
	w.AddMapping("POST", "/path1", NewStubHTTPResponse(200, "test body"))
	status, _, reader, _, err := w.Handle(
		context.Background(),
		"/path1",
		"POST",
		map[string][]string{"key": {"value"}},
		make(map[string]string),
		nil)
	require.NoError(t, err)
	require.Equal(t, 200, status)
	b, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "test body", string(b))
}

func Test_ShouldPostStubClientWithErrorResponse(t *testing.T) {
	w := NewStubHTTPClient()
	w.AddMapping("POST", "/path1", NewStubHTTPResponseError(500, 1, fmt.Errorf("test error")))
	_, _, _, _, err := w.Handle(
		context.Background(),
		"/path1",
		"POST",
		map[string][]string{"key": {"value"}},
		make(map[string]string),
		nil)
	require.Error(t, err)
}
