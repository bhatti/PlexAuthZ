package web

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// ********************************* STUB METHODS For Web server ***********************************
type stubWebServer struct {
}

type stubContext struct {
	request  *http.Request
	response *echo.Response
	Result   any
	Params   map[string]string
	Context  map[string]any
}

func (c *stubContext) SetResponse(r *echo.Response) {
	c.response = r
}

func (c *stubContext) SetLogger(_ echo.Logger) {
}

// NewStubContext - creates stubbed server
func NewStubContext(req *http.Request) *stubContext { //nolint
	return &stubContext{
		request:  req,
		response: echo.NewResponse(&StubResponseWriter{}, nil),
		Params:   make(map[string]string),
		Context:  make(map[string]any),
	}
}

// NewStubWebServer creates stubbed web server
func NewStubWebServer() Server {
	return &stubWebServer{}
}

// AddMiddleware middleware
func (w *stubWebServer) AddMiddleware(_ echo.MiddlewareFunc) {
}

func (w *stubWebServer) GET(string, HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

func (w *stubWebServer) POST(string, HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

func (w *stubWebServer) PUT(string, HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

func (w *stubWebServer) DELETE(string, HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

// CONNECT calls HTTP CONNECT method
func (w *stubWebServer) CONNECT(_ string, _ HandlerFunc, _ ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

// HEAD calls HTTP HEAD method
func (w *stubWebServer) HEAD(_ string, _ HandlerFunc, _ ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

// OPTIONS calls HTTP OPTIONS method
func (w *stubWebServer) OPTIONS(_ string, _ HandlerFunc, _ ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

// PATCH calls HTTP PATCH method
func (w *stubWebServer) PATCH(_ string, _ HandlerFunc, _ ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

// TRACE calls HTTP TRACE method
func (w *stubWebServer) TRACE(_ string, _ HandlerFunc, _ ...echo.MiddlewareFunc) *echo.Route {
	return &echo.Route{}
}

func (w *stubWebServer) Start(string) error {
	return nil
}

func (w *stubWebServer) Addr() string {
	return ""
}

func (w *stubWebServer) Stop() error {
	return nil
}

func (w *stubWebServer) Static(string, string) {
}

func (w *stubWebServer) Embed(embed.FS, string, string) {
}

// ********************************* STUB METHODS For Echo Context ***********************************

// Request getter
func (c *stubContext) Request() *http.Request {
	return c.request
}

// SetRequest setter
func (c *stubContext) SetRequest(r *http.Request) {
	c.request = r
}

// Response getter
func (c *stubContext) Response() *echo.Response {
	return c.response
}

// IsTLS getter
func (c *stubContext) IsTLS() bool {
	return true
}

// IsWebSocket getter
func (c *stubContext) IsWebSocket() bool {
	return false
}

// Scheme getter
func (c *stubContext) Scheme() string {
	return "https"
}

func (c *stubContext) RealIP() string {
	return "127.0.0.1"
}

func (c *stubContext) Path() string {
	return ""
}

func (c *stubContext) SetPath(string) {
}

func (c *stubContext) Param(name string) string {
	return c.Params[name]
}

func (c *stubContext) ParamNames() []string {
	return make([]string, 0)
}

func (c *stubContext) SetParamNames(...string) {
}

func (c *stubContext) ParamValues() []string {
	return make([]string, 0)
}

func (c *stubContext) SetParamValues(...string) {
}

func (c *stubContext) QueryParam(name string) string {
	return c.Params[name]
}

func (c *stubContext) QueryParams() (res url.Values) {
	res = make(url.Values)
	for k, v := range c.Params {
		res[k] = []string{v}
	}
	if c.request.URL != nil {
		for k, v := range c.request.URL.Query() {
			res[k] = v
		}
	}
	return
}

func (c *stubContext) QueryString() string {
	return ""
}

func (c *stubContext) FormValue(name string) string {
	return c.Params[name]
}

func (c *stubContext) FormParams() (url.Values, error) {
	return c.request.Form, nil
}

func (c *stubContext) FormFile(string) (*multipart.FileHeader, error) {
	return nil, nil
}

func (c *stubContext) MultipartForm() (*multipart.Form, error) {
	return nil, nil
}

func (c *stubContext) Cookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

func (c *stubContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response(), cookie)
}

func (c *stubContext) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

func (c *stubContext) Get(name string) any {
	return c.Context[name]
}

func (c *stubContext) Set(name string, value any) {
	c.Context[name] = value
}

func (c *stubContext) Bind(any) error {
	return nil
}

func (c *stubContext) Validate(any) error {
	return nil
}

func (c *stubContext) Render(code int, t string, r any) (err error) {
	c.Result = r
	if code >= 300 {
		return fmt.Errorf("%d - %s", code, t)
	}
	return nil
}

func (c *stubContext) HTML(code int, t string) (err error) {
	c.Result = t
	if code >= 300 {
		return fmt.Errorf("%d - %s", code, t)
	}
	return nil
}

func (c *stubContext) HTMLBlob(code int, b []byte) (err error) {
	c.Result = b
	if code >= 300 {
		return fmt.Errorf("%d - %d", code, len(b))
	}
	return nil
}

func (c *stubContext) String(code int, s string) (err error) {
	c.Result = s
	if code >= 300 {
		return fmt.Errorf("%d - %s", code, s)
	}
	return nil
}

func (c *stubContext) JSON(code int, j any) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %v", code, j)
	}
	return nil
}

func (c *stubContext) JSONPretty(code int, j any, _ string) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %v", code, j)
	}
	return nil
}

func (c *stubContext) JSONBlob(code int, b []byte) (err error) {
	c.Result = b
	if code >= 300 {
		return fmt.Errorf("%d - %d", code, len(b))
	}
	return nil
}

func (c *stubContext) JSONP(code int, j string, _ any) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %s", code, j)
	}
	return nil
}

func (c *stubContext) JSONPBlob(code int, j string, _ []byte) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %s", code, j)
	}
	return nil
}

func (c *stubContext) XML(code int, j any) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %v", code, j)
	}
	return nil
}

func (c *stubContext) XMLPretty(code int, j any, _ string) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %v", code, j)
	}
	return nil
}

func (c *stubContext) XMLBlob(code int, j []byte) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %v", code, j)
	}
	return nil
}

func (c *stubContext) Blob(code int, _ string, j []byte) (err error) {
	c.Result = j
	if code >= 300 {
		return fmt.Errorf("%d - %v", code, j)
	}
	return nil
}

func (c *stubContext) Stream(code int, _ string, r io.Reader) (err error) {
	c.Result, _ = io.ReadAll(r)
	if code >= 300 {
		return fmt.Errorf("%d - %v", code, r)
	}
	return nil
}

func (c *stubContext) File(string) (err error) {
	return
}

func (c *stubContext) Attachment(string, string) error {
	return nil
}

func (c *stubContext) Inline(string, string) error {
	return nil
}

func (c *stubContext) NoContent(code int) error {
	if code >= 300 {
		return fmt.Errorf("%d", code)
	}
	return nil
}

func (c *stubContext) Redirect(code int, s string) error {
	if code >= 300 {
		return fmt.Errorf("%d u %s", code, s)
	}
	return nil
}

func (c *stubContext) Error(err error) {
	c.Result = err
}

func (c *stubContext) Echo() *echo.Echo {
	return nil
}

func (c *stubContext) Handler() echo.HandlerFunc {
	return nil
}

func (c *stubContext) SetHandler(echo.HandlerFunc) {
}

func (c *stubContext) Logger() echo.Logger {
	return nil
}

func (c *stubContext) Reset(*http.Request, http.ResponseWriter) {
}
