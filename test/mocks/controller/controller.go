package controller

import (
    "net/http"
    "strings"

    "github.com/jamiefdhurst/journal/pkg/session"
)

// MockController Mock the controller interface
type MockController struct {
    HasRun  bool
    session *session.Session
}

// Init Mock the init method
func (m *MockController) Init(app interface{}, params []string, request *http.Request) {
    m.session = session.NewSession()
}

// Run Mock the run method
func (m *MockController) Run(response http.ResponseWriter, request *http.Request) {
    m.HasRun = true
}

func (m *MockController) Container() interface{} {
    var r interface{}
    return r
}

func (m *MockController) Host() string {
    return "foobar.com"
}

func (m *MockController) Params() []string {
    return []string{}
}

func (m *MockController) SaveSession(w http.ResponseWriter) {}

func (m *MockController) Session() *session.Session {
    return m.session
}

// MockResponse Mock for http.ResponseWriter
type MockResponse struct {
    Content    string
    Headers    http.Header
    StatusCode int
}

// Header Return Headers map
func (m *MockResponse) Header() http.Header {
    return m.Headers
}

// Reset Reset the struct
func (m *MockResponse) Reset() {
    m.Content = ""
    m.Headers = make(http.Header)
    m.StatusCode = 200
}

// Write Write the response
func (m *MockResponse) Write(b []byte) (int, error) {
    m.Content = strings.Join([]string{m.Content, string(b[:])}, "")
    return len(b), nil
}

// WriteHeader Write the status code
func (m *MockResponse) WriteHeader(statusCode int) {
    m.StatusCode = statusCode
}

// NewMockResponse Make a mock response
func NewMockResponse() *MockResponse {
    m := &MockResponse{}
    m.Reset()
    return m
}
