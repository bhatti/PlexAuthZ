package web

import (
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// APIContext interface
type APIContext interface { //nolint
	// Path Request path
	Path() string

	// Request returns `*http.Request`.
	Request() *http.Request

	// Response returns `*Response`.
	Response() *echo.Response

	// Param returns path parameter by name.
	Param(name string) string

	// QueryParams returns the query parameters as `url.Values`.
	QueryParams() url.Values

	// QueryParam returns the query param for the provided name.
	QueryParam(name string) string

	// FormParams returns the form parameters as `url.Values`.
	FormParams() (url.Values, error)

	// FormValue returns the form field value for the provided name.
	FormValue(name string) string

	// Cookie returns the named cookie provided in the request.
	Cookie(name string) (*http.Cookie, error)

	// SetCookie adds a `Set-Cookie` header in HTTP response.
	SetCookie(cookie *http.Cookie)

	// Get retrieves data from the context.
	Get(key string) any

	// Set saves data in the context.
	Set(key string, val any)

	// Render renders a template with data and sends a text/html response with status
	// code. Renderer must be registered using `Echo.Renderer`.
	Render(code int, name string, data any) error

	// String sends a string response with status code.
	String(code int, s string) error

	// JSON sends a JSON response with status code.
	JSON(code int, i any) error

	// MultipartForm returns the multipart form.
	MultipartForm() (*multipart.Form, error)

	// Redirect redirects the request to a provided URL with status code.
	Redirect(code int, url string) error

	// NoContent sends a response with nobody and a status code.
	NoContent(code int) error

	// Blob sends a blob response with status code and content type.
	Blob(code int, contentType string, b []byte) error

	// Stream sends a streaming response with status code and content type.
	Stream(code int, contentType string, r io.Reader) error
	// Attachment sends a response as attachment, prompting client to save the
	// file.
	Attachment(file string, name string) error
}
