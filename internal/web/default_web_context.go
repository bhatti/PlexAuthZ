package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// DefaultAPIContext struct for defining APIContext using http request/response
type DefaultAPIContext struct {
	ctx       context.Context
	request   *http.Request
	response  *http.Response
	eResponse *echo.Response
	params    map[string]string
}

// NewDefaultAPIContext constructor
func NewDefaultAPIContext(
	ctx context.Context,
	request *http.Request,
	response *http.Response,
	eResponse *echo.Response,
	params map[string]string) APIContext {
	return &DefaultAPIContext{
		ctx:       ctx,
		request:   request,
		response:  response,
		eResponse: eResponse,
		params:    params,
	}
}

// Path Request path
func (d *DefaultAPIContext) Path() string {
	return d.request.URL.Path
}

// Host Request host
func (d *DefaultAPIContext) Host() string {
	return d.request.Host
}

// Request returns `*http.Request`.
func (d *DefaultAPIContext) Request() *http.Request {
	return d.request
}

// Response returns `*Response`.
func (d *DefaultAPIContext) Response() *echo.Response {
	return d.eResponse
}

// Param returns path parameter by name.
func (d *DefaultAPIContext) Param(name string) string {
	return d.params[name]
}

// QueryParams returns the query parameters as `url.Values`.
func (d *DefaultAPIContext) QueryParams() url.Values {
	return d.request.URL.Query()
}

// QueryParam returns the query param for the provided name.
func (d *DefaultAPIContext) QueryParam(name string) string {
	return d.request.URL.Query().Get(name)
}

// FormParams returns the form parameters as `url.Values`.
func (d *DefaultAPIContext) FormParams() (url.Values, error) {
	return d.request.Form, nil
}

// FormValue returns the form field value for the provided name.
func (d *DefaultAPIContext) FormValue(name string) string {
	return d.request.Form.Get(name)
}

// Cookie returns the named cookie provided in the request.
func (d *DefaultAPIContext) Cookie(name string) (*http.Cookie, error) {
	return d.request.Cookie(name)
}

// SetCookie adds a `Set-Cookie` header in HTTP response.
func (d *DefaultAPIContext) SetCookie(cookie *http.Cookie) {
	d.response.Header.Add("Set-Cookie", cookie.String())
}

// Get retrieves data from the context.
func (d *DefaultAPIContext) Get(key string) any {
	return d.ctx.Value(key)
}

// Set saves data in the context.
func (d *DefaultAPIContext) Set(key string, val any) {
	d.ctx = context.WithValue(d.ctx, key, val)
}

// Render renders a template with data and sends a text/html response with status
// code. Renderer must be registered using `Echo.Renderer`.
func (d *DefaultAPIContext) Render(_ int, _ string, _ any) error {
	return fmt.Errorf("render not supported")
}

// String sends a string response with status code.
func (d *DefaultAPIContext) String(code int, s string) error {
	d.response.StatusCode = code
	d.response.Body = io.NopCloser(bytes.NewReader([]byte(s)))
	return nil
}

// JSON sends a JSON response with status code.
func (d *DefaultAPIContext) JSON(code int, i any) error {
	d.response.StatusCode = code
	j, err := json.Marshal(i)
	if err != nil {
		return err
	}
	d.response.Body = io.NopCloser(bytes.NewReader(j))
	return nil
}

// MultipartForm returns the multipart form.
func (d *DefaultAPIContext) MultipartForm() (*multipart.Form, error) {
	return nil, fmt.Errorf("multipartForm not supported")
}

// Redirect redirects the request to a provided URL with status code.
func (d *DefaultAPIContext) Redirect(code int, url string) error {
	d.response.StatusCode = code
	d.response.Header.Add("Location", url)
	return nil
}

// NoContent sends a response with nobody and a status code.
func (d *DefaultAPIContext) NoContent(code int) error {
	d.response.StatusCode = code
	return nil
}

// Blob sends a blob response with status code and content type.
func (d *DefaultAPIContext) Blob(code int, contentType string, b []byte) error {
	d.response.StatusCode = code
	d.response.Header.Add("Content-Type", contentType)
	d.response.Body = io.NopCloser(bytes.NewReader(b))
	return nil
}

// Stream sends a streaming response with status code and content type.
func (d *DefaultAPIContext) Stream(code int, contentType string, r io.Reader) (err error) {
	d.response.StatusCode = code
	d.response.Header.Add("Content-Type", contentType)
	d.response.Body = io.NopCloser(r)
	return nil
}

// Attachment sends a response as attachment, prompting client to save the
// file.
func (d *DefaultAPIContext) Attachment(_ string, _ string) error {
	return fmt.Errorf("attachment not supported")
}
