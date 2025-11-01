package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDefaultStore(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		config      CookieConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid 32-byte key",
			key:  "12345678901234567890123456789012",
			config: CookieConfig{
				Name:     "test-session",
				Domain:   "example.com",
				MaxAge:   3600,
				Secure:   true,
				HTTPOnly: true,
			},
			expectError: false,
		},
		{
			name: "Key too short",
			key:  "tooshort",
			config: CookieConfig{
				Name: "test-session",
			},
			expectError: true,
			errorMsg:    "session key must be exactly 32 bytes",
		},
		{
			name: "Key too long",
			key:  "123456789012345678901234567890123",
			config: CookieConfig{
				Name: "test-session",
			},
			expectError: true,
			errorMsg:    "session key must be exactly 32 bytes",
		},
		{
			name: "Invalid characters in key",
			key:  "123456789012345678901234\x00\x01\x02\x03\x04\x05\x06\x07",
			config: CookieConfig{
				Name: "test-session",
			},
			expectError: true,
			errorMsg:    "session key must contain only printable ASCII characters",
		},
		{
			name: "Default cookie name when empty",
			key:  "12345678901234567890123456789012",
			config: CookieConfig{
				Name: "",
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, err := NewDefaultStore(test.key, test.config)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != test.errorMsg {
					t.Errorf("Expected error %q, got %q", test.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if store == nil {
					t.Errorf("Expected store to be created but got nil")
				}
				if test.config.Name == "" && store.name != "journal-session" {
					t.Errorf("Expected default name 'journal-session', got %q", store.name)
				}
				if test.config.Name != "" && store.name != test.config.Name {
					t.Errorf("Expected name %q, got %q", test.config.Name, store.name)
				}
			}
		})
	}
}

func TestEncryptDecryptCycle(t *testing.T) {
	key := "12345678901234567890123456789012"
	config := CookieConfig{
		Name:     "test-session",
		Domain:   "",
		MaxAge:   3600,
		Secure:   false,
		HTTPOnly: true,
	}

	store, err := NewDefaultStore(key, config)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	testData := map[string]interface{}{
		"user_id": "12345",
		"name":    "Test User",
		"count":   42,
		"active":  true,
	}

	encrypted, err := store.encrypt(testData)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if encrypted == "" {
		t.Errorf("Encrypted string should not be empty")
	}

	var decrypted map[string]interface{}
	err = store.decrypt(encrypted, &decrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted["user_id"] != testData["user_id"] {
		t.Errorf("Expected user_id %v, got %v", testData["user_id"], decrypted["user_id"])
	}
	if decrypted["name"] != testData["name"] {
		t.Errorf("Expected name %v, got %v", testData["name"], decrypted["name"])
	}
}

func TestCookieConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config CookieConfig
	}{
		{
			name: "Secure cookie with HTTPOnly",
			config: CookieConfig{
				Name:     "secure-session",
				Domain:   "example.com",
				MaxAge:   7200,
				Secure:   true,
				HTTPOnly: true,
			},
		},
		{
			name: "Non-secure cookie without HTTPOnly",
			config: CookieConfig{
				Name:     "insecure-session",
				Domain:   "",
				MaxAge:   3600,
				Secure:   false,
				HTTPOnly: false,
			},
		},
		{
			name: "Custom domain cookie",
			config: CookieConfig{
				Name:     "domain-session",
				Domain:   "example.com",
				MaxAge:   1800,
				Secure:   true,
				HTTPOnly: true,
			},
		},
		{
			name: "Long expiry cookie",
			config: CookieConfig{
				Name:     "long-session",
				Domain:   "",
				MaxAge:   2592000,
				Secure:   false,
				HTTPOnly: true,
			},
		},
	}

	key := "12345678901234567890123456789012"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, err := NewDefaultStore(key, test.config)
			if err != nil {
				t.Fatalf("Failed to create store: %v", err)
			}

			session := NewSession()
			session.Set("test", "value")
			store.cachedSession = session

			w := httptest.NewRecorder()
			err = store.Save(w)
			if err != nil {
				t.Fatalf("Failed to save session: %v", err)
			}

			cookies := w.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("Expected 1 cookie, got %d", len(cookies))
			}

			cookie := cookies[0]

			if cookie.Name != test.config.Name {
				t.Errorf("Expected cookie name %q, got %q", test.config.Name, cookie.Name)
			}
			if cookie.Domain != test.config.Domain {
				t.Errorf("Expected cookie domain %q, got %q", test.config.Domain, cookie.Domain)
			}
			if cookie.MaxAge != test.config.MaxAge {
				t.Errorf("Expected cookie MaxAge %d, got %d", test.config.MaxAge, cookie.MaxAge)
			}
			if cookie.Secure != test.config.Secure {
				t.Errorf("Expected cookie Secure %v, got %v", test.config.Secure, cookie.Secure)
			}
			if cookie.HttpOnly != test.config.HTTPOnly {
				t.Errorf("Expected cookie HttpOnly %v, got %v", test.config.HTTPOnly, cookie.HttpOnly)
			}
			if cookie.Path != "/" {
				t.Errorf("Expected cookie Path '/', got %q", cookie.Path)
			}
			if cookie.SameSite != http.SameSiteStrictMode {
				t.Errorf("Expected cookie SameSite Strict, got %v", cookie.SameSite)
			}
		})
	}
}

func TestGetSession(t *testing.T) {
	key := "12345678901234567890123456789012"
	config := CookieConfig{
		Name:     "test-session",
		Domain:   "",
		MaxAge:   3600,
		Secure:   false,
		HTTPOnly: true,
	}

	store, err := NewDefaultStore(key, config)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	t.Run("Get session without cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		session, err := store.Get(req)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if session == nil {
			t.Errorf("Expected session to be created")
		}
	})

	t.Run("Get session with valid cookie", func(t *testing.T) {
		session := NewSession()
		session.Set("user", "testuser")
		store.cachedSession = session

		w := httptest.NewRecorder()
		err := store.Save(w)
		if err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatalf("Expected 1 cookie, got %d", len(cookies))
		}

		newStore, err := NewDefaultStore(key, config)
		if err != nil {
			t.Fatalf("Failed to create new store: %v", err)
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(cookies[0])

		retrievedSession, err := newStore.Get(req)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if retrievedSession == nil {
			t.Fatalf("Expected session to be retrieved")
		}

		user := retrievedSession.Get("user")
		if user == nil {
			t.Errorf("Expected 'user' key to exist in session")
		}
		if user != "testuser" {
			t.Errorf("Expected user 'testuser', got %v", user)
		}
	})
}

func TestSessionCaching(t *testing.T) {
	key := "12345678901234567890123456789012"
	config := CookieConfig{
		Name:     "test-session",
		Domain:   "",
		MaxAge:   3600,
		Secure:   false,
		HTTPOnly: true,
	}

	store, err := NewDefaultStore(key, config)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	session1, err := store.Get(req)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	session2, err := store.Get(req)
	if err != nil {
		t.Fatalf("Failed to get session second time: %v", err)
	}

	if session1 != session2 {
		t.Errorf("Expected same session instance to be returned (cached)")
	}
}
