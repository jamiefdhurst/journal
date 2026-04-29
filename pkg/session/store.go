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

// Store defines the interface for session storage implementations
type Store interface {
    Get(r *http.Request) (*Session, error)
    Save(w http.ResponseWriter) error
}

const defaultName string = "journal-session"

// CookieConfig defines the configuration for session cookies
type CookieConfig struct {
    Name     string
    Domain   string
    MaxAge   int
    Secure   bool
    HTTPOnly bool
}

// DefaultStore implements Store using encrypted cookies for session storage
type DefaultStore struct {
    cachedSession *Session
    key           []byte
    name          string
    config        CookieConfig
}

// NewDefaultStore creates a new DefaultStore with the given encryption key and cookie configuration.
// The key must be exactly 32 bytes (for AES-256) and contain only printable ASCII characters.
func NewDefaultStore(key string, config CookieConfig) (*DefaultStore, error) {
    if len(key) != 32 {
        return nil, errors.New("session key must be exactly 32 bytes")
    }

    for i := 0; i < len(key); i++ {
        if key[i] < 32 || key[i] > 126 {
            return nil, errors.New("session key must contain only printable ASCII characters")
        }
    }

    name := config.Name
    if name == "" {
        name = defaultName
    }

    return &DefaultStore{
        key:    []byte(key),
        name:   name,
        config: config,
    }, nil
}

// Get retrieves the session from the request cookie, decrypting and deserializing it.
// If no session exists, a new empty session is created.
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
        } else {
            s.cachedSession = NewSession()
        }
    }

    return s.cachedSession, err
}

// Save encrypts and serializes the session, writing it to a cookie in the response.
func (s *DefaultStore) Save(w http.ResponseWriter) error {
    encrypted, err := s.encrypt(s.cachedSession.Values)
    if err != nil {
        return err
    }

    http.SetCookie(w, &http.Cookie{
        Name:     s.name,
        Value:    encrypted,
        Path:     "/",
        Domain:   s.config.Domain,
        MaxAge:   s.config.MaxAge,
        Secure:   s.config.Secure,
        SameSite: http.SameSiteStrictMode,
        HttpOnly: s.config.HTTPOnly,
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
    asBytes, err := base64.URLEncoding.DecodeString(encrypted)
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
