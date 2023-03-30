package controller

import (
	"net/http"
	"strings"
)

// MockController Mock the controller interface
type MockController struct {
	HasRun bool
}

// Init Mock the init method
func (m *MockController) Init(app interface{}, params []string, request *http.Request) {}

// Run Mock the run method
func (m *MockController) Run(response http.ResponseWriter, request *http.Request) {
	m.HasRun = true
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
	m.Headers = make(http.Header)
	return m
}
