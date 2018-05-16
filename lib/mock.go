package lib

// MockServer Mock the server interface
type MockServer struct{}

// ListenAndServe Dummy method
func (m MockServer) ListenAndServe() error {
	return nil
}
