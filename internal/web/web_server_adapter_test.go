package web

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func adapterHandler(c APIContext) error {
	res := make(map[string]string)
	res["url-path"] = c.Request().URL.Path
	res["http-path"] = c.Param("path")
	res["http-method"] = c.Param("method")
	res["http-name"] = c.Param("name")
	for k, v := range c.Request().URL.Query() {
		res["query-"+k] = v[0]
	}
	return c.JSON(200, res)
}

func Test_ShouldNotInvokeHTTPRequestWithUnknownPath(t *testing.T) {
	// GIVEN a server adapter
	adapter := NewWebServerAdapter()
	// WHEN using unknown path
	u, err := url.Parse("http://localhost:8080/abc")
	require.NoError(t, err)
	req := &http.Request{
		URL:    u,
		Method: "POST",
		Header: http.Header{"X1": []string{"val1"}, "Content-Type": []string{"json"}},
	}
	// THEN it should fail to invoke
	res, err := adapter.Invoke(req)
	require.Error(t, err)
	require.Nil(t, res)
}

func Test_ShouldNotInvokeHTTPRequestWithUnknownMethod(t *testing.T) {
	// GIVEN a server adapter
	adapter := NewWebServerAdapter()
	// WHEN using unknown path
	u, err := url.Parse("http://localhost:8080/abc")
	require.NoError(t, err)
	req := &http.Request{
		URL:    u,
		Method: "XXXX",
		Header: http.Header{"X1": []string{"val1"}, "Content-Type": []string{"json"}},
	}
	// THEN it should fail to invoke
	res, err := adapter.Invoke(req)
	require.Error(t, err)
	require.Nil(t, res)
}
