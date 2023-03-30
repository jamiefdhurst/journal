package session

const flashKey = "_flash"

type Session struct {
	ID     string
	Values map[string]interface{}
}

func NewSession() *Session {
	return &Session{
		Values: make(map[string]interface{}),
	}
}

func (s *Session) GetFlash() []interface{} {
	var flashes []interface{}

	if v, ok := s.Values[flashKey]; ok {
		delete(s.Values, flashKey)
		flashes = v.([]interface{})
	}

	return flashes
}

func (s *Session) AddFlash(value interface{}) {
	var flashes []interface{}
	if v, ok := s.Values[flashKey]; ok {
		flashes = v.([]interface{})
	}
	s.Values[flashKey] = append(flashes, value)
}
