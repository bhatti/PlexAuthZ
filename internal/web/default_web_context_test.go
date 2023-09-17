package web

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func Test_ShouldGetDefaultAPIContextPath(t *testing.T) {
	request, ctx := newTestWebContext()
	require.Equal(t, "/path", ctx.Path())
	require.Equal(t, request, ctx.Request())
}

func Test_ShouldNotGetResponseForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	require.Nil(t, ctx.Response())
}

func Test_ShouldGetParamForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	require.Equal(t, "one", ctx.Param("key"))
}

func Test_ShouldGetQueryParamForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	require.True(t, len(ctx.QueryParams()) > 0)
	require.Equal(t, "a", ctx.QueryParam("x"))
}

func Test_ShouldNotGetFormParamForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	require.Equal(t, "", ctx.FormValue("abc"))
	form, err := ctx.FormParams()
	require.NoError(t, err)
	require.Equal(t, 0, len(form))
}

func Test_ShouldGetSetCookieForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	cookie, err := ctx.Cookie("test")
	require.Error(t, err)
	require.Nil(t, cookie)
	ctx.SetCookie(&http.Cookie{Name: "name", Value: "value"})
}

func Test_ShouldGetSetValueForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	ctx.Set("xyz", "123")
	require.Equal(t, "123", ctx.Get("xyz"))
}

func Test_ShouldNotRenderAttachMultipartFormForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	require.Error(t, ctx.Render(1, "", ""))
	require.Error(t, ctx.Attachment("file", ""))
	_, err := ctx.MultipartForm()
	require.Error(t, err)
}

func Test_ShouldReturnResponseForDefaultAPIContextPath(t *testing.T) {
	_, ctx := newTestWebContext()
	require.NoError(t, ctx.Blob(200, "", []byte("test")))
	require.NoError(t, ctx.Stream(200, "", bytes.NewReader([]byte("test"))))
	require.NoError(t, ctx.String(200, "test"))
	require.NoError(t, ctx.JSON(200, "test"))
	require.NoError(t, ctx.Redirect(301, "test"))
	require.NoError(t, ctx.NoContent(200))
}

func newTestWebContext() (*http.Request, APIContext) {
	request := &http.Request{}
	request.Header = make(http.Header)
	request.Body = io.NopCloser(bytes.NewReader([]byte{}))
	request.URL, _ = url.Parse("http://localhost/path?x=a")
	response := &http.Response{}
	response.Request = request
	response.Header = make(http.Header)
	response.Body = io.NopCloser(bytes.NewReader([]byte{}))
	ctx := NewDefaultAPIContext(context.Background(), request, response, nil, map[string]string{"key": "one"})
	return request, ctx
}
