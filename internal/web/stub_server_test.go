package web

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_ShouldNewStubWebResponse(t *testing.T) {
	m := NewStubHTTPResponse(200, "test")
	require.NotNil(t, m)
}

func Test_ShouldNewStubUtils(t *testing.T) {
	w := NewStubHTTPClient()
	require.NotNil(t, w)
}

func Test_ShouldNewStubEchoContext(t *testing.T) {
	w := NewStubContext(&http.Request{})
	require.NotNil(t, w)
}

func Test_ShouldNewStubWebServer(t *testing.T) {
	w := NewStubWebServer()
	require.NotNil(t, w)
}
