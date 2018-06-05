package router

// MockServer Mock the server interface
type MockServer struct {
	Listening bool
}

// ListenAndServe Dummy method
func (m *MockServer) ListenAndServe() error {
	m.Listening = true
	return nil
}
