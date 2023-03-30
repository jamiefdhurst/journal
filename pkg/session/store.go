package session

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io"
	"net/http"
)

type Store interface {
	Get(r *http.Request) (*Session, error)
	Save(w http.ResponseWriter) error
}

const defaultName string = "journal-session"

type DefaultStore struct {
	cachedSession *Session
	key           []byte
	name          string
}

func NewDefaultStore(key string) *DefaultStore {
	return &DefaultStore{
		key:  []byte(key),
		name: defaultName,
	}
}

func (s *DefaultStore) Get(r *http.Request) (*Session, error) {
	var err error
	if s.cachedSession == nil {
		session := NewSession()
		c, err := r.Cookie(s.name)
		if err == nil {
			err = s.decrypt(c.Value, &session.Values)
		} else {
			err = nil
		}
		if err == nil {
			s.cachedSession = session
		}
	}

	return s.cachedSession, err
}

func (s *DefaultStore) Save(w http.ResponseWriter) error {
	encrypted, err := s.encrypt(s.cachedSession.Values)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.name,
		Value:    encrypted,
		Path:     "/",
		Domain:   "",
		MaxAge:   86400 * 30,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: false,
	})

	return nil
}

func (s *DefaultStore) decrypt(encrypted string, output interface{}) error {
	c, err := aes.NewCipher(s.key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	asBytes, _ := base64.URLEncoding.DecodeString(encrypted)
	if err != nil || len(asBytes) < nonceSize {
		return errors.New("string length too short")
	}

	nonce, asBytes := asBytes[:nonceSize], asBytes[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, asBytes, nil)
	if err != nil {
		return err
	}

	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	dec := gob.NewDecoder(bytes.NewBuffer(decrypted))
	return dec.Decode(output)
}

func (s *DefaultStore) encrypt(input interface{}) (string, error) {
	var buf bytes.Buffer
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(input); err != nil {
		return "", err
	}

	c, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, buf.Bytes(), nil)), nil
}
