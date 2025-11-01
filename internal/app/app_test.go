package app

import (
	"os"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	config := DefaultConfiguration()

	if config.ArticlesPerPage != 20 {
		t.Errorf("Expected ArticlesPerPage 20, got %d", config.ArticlesPerPage)
	}
	if config.Port != "3000" {
		t.Errorf("Expected Port '3000', got %q", config.Port)
	}
	if config.SessionName != "journal-session" {
		t.Errorf("Expected SessionName 'journal-session', got %q", config.SessionName)
	}
	if config.CookieMaxAge != 2592000 {
		t.Errorf("Expected CookieMaxAge 2592000, got %d", config.CookieMaxAge)
	}
	if config.CookieHTTPOnly != true {
		t.Errorf("Expected CookieHTTPOnly true, got %v", config.CookieHTTPOnly)
	}
	if config.CookieSecure != false {
		t.Errorf("Expected CookieSecure false, got %v", config.CookieSecure)
	}
	if config.SessionKey != "" {
		t.Errorf("Expected SessionKey to be empty by default, got %q", config.SessionKey)
	}
}

func TestApplyEnvConfiguration_SessionKey(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		expectWarning bool
		expectKey     bool
	}{
		{
			name:          "Valid 32-byte key",
			envValue:      "12345678901234567890123456789012",
			expectWarning: false,
			expectKey:     true,
		},
		{
			name:          "Key too short generates auto key",
			envValue:      "tooshort",
			expectWarning: true,
			expectKey:     true,
		},
		{
			name:          "Key too long generates auto key",
			envValue:      "123456789012345678901234567890123",
			expectWarning: true,
			expectKey:     true,
		},
		{
			name:          "Empty key generates auto key",
			envValue:      "",
			expectWarning: true,
			expectKey:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("J_SESSION_KEY", test.envValue)
			defer os.Unsetenv("J_SESSION_KEY")

			config := DefaultConfiguration()
			ApplyEnvConfiguration(&config)

			if test.expectKey && config.SessionKey == "" {
				t.Errorf("Expected session key to be set")
			}
			if test.expectKey && len(config.SessionKey) != 32 {
				t.Errorf("Expected session key length 32, got %d", len(config.SessionKey))
			}
			if test.envValue != "" && len(test.envValue) == 32 && config.SessionKey != test.envValue {
				t.Errorf("Expected session key %q, got %q", test.envValue, config.SessionKey)
			}
		})
	}
}

func TestApplyEnvConfiguration_SessionName(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "Custom session name",
			envValue: "custom-session",
			expected: "custom-session",
		},
		{
			name:     "Empty uses default",
			envValue: "",
			expected: "journal-session",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv("J_SESSION_NAME", test.envValue)
				defer os.Unsetenv("J_SESSION_NAME")
			}

			config := DefaultConfiguration()
			ApplyEnvConfiguration(&config)

			if config.SessionName != test.expected {
				t.Errorf("Expected SessionName %q, got %q", test.expected, config.SessionName)
			}
		})
	}
}

func TestApplyEnvConfiguration_CookieDomain(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "Custom domain",
			envValue: ".example.com",
			expected: ".example.com",
		},
		{
			name:     "Specific subdomain",
			envValue: "app.example.com",
			expected: "app.example.com",
		},
		{
			name:     "Empty uses default",
			envValue: "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv("J_COOKIE_DOMAIN", test.envValue)
				defer os.Unsetenv("J_COOKIE_DOMAIN")
			}

			config := DefaultConfiguration()
			ApplyEnvConfiguration(&config)

			if config.CookieDomain != test.expected {
				t.Errorf("Expected CookieDomain %q, got %q", test.expected, config.CookieDomain)
			}
		})
	}
}

func TestApplyEnvConfiguration_CookieMaxAge(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected int
	}{
		{
			name:     "Custom max age",
			envValue: "7200",
			expected: 7200,
		},
		{
			name:     "One week",
			envValue: "604800",
			expected: 604800,
		},
		{
			name:     "Invalid uses default",
			envValue: "invalid",
			expected: 2592000,
		},
		{
			name:     "Empty uses default",
			envValue: "",
			expected: 2592000,
		},
		{
			name:     "Zero uses default",
			envValue: "0",
			expected: 2592000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv("J_COOKIE_MAX_AGE", test.envValue)
				defer os.Unsetenv("J_COOKIE_MAX_AGE")
			}

			config := DefaultConfiguration()
			ApplyEnvConfiguration(&config)

			if config.CookieMaxAge != test.expected {
				t.Errorf("Expected CookieMaxAge %d, got %d", test.expected, config.CookieMaxAge)
			}
		})
	}
}

func TestApplyEnvConfiguration_CookieHTTPOnly(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "Disabled with 0",
			envValue: "0",
			expected: false,
		},
		{
			name:     "Disabled with false",
			envValue: "false",
			expected: false,
		},
		{
			name:     "Enabled with 1",
			envValue: "1",
			expected: true,
		},
		{
			name:     "Enabled with true",
			envValue: "true",
			expected: true,
		},
		{
			name:     "Default is enabled",
			envValue: "",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv("J_COOKIE_HTTPONLY", test.envValue)
				defer os.Unsetenv("J_COOKIE_HTTPONLY")
			}

			config := DefaultConfiguration()
			ApplyEnvConfiguration(&config)

			if config.CookieHTTPOnly != test.expected {
				t.Errorf("Expected CookieHTTPOnly %v, got %v", test.expected, config.CookieHTTPOnly)
			}
		})
	}
}

func TestApplyEnvConfiguration_CookieSecure(t *testing.T) {
	tests := []struct {
		name        string
		sslCert     string
		sslKey      string
		expected    bool
		description string
	}{
		{
			name:        "Secure when SSL cert is set",
			sslCert:     "/path/to/cert.pem",
			sslKey:      "/path/to/key.pem",
			expected:    true,
			description: "Cookie should be secure when SSL is enabled",
		},
		{
			name:        "Not secure when SSL cert is empty",
			sslCert:     "",
			sslKey:      "",
			expected:    false,
			description: "Cookie should not be secure when SSL is not enabled",
		},
		{
			name:        "Secure even without key if cert is set",
			sslCert:     "/path/to/cert.pem",
			sslKey:      "",
			expected:    true,
			description: "Cookie secure flag follows cert presence",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.sslCert != "" {
				os.Setenv("J_SSL_CERT", test.sslCert)
				defer os.Unsetenv("J_SSL_CERT")
			}
			if test.sslKey != "" {
				os.Setenv("J_SSL_KEY", test.sslKey)
				defer os.Unsetenv("J_SSL_KEY")
			}

			config := DefaultConfiguration()
			ApplyEnvConfiguration(&config)

			if config.CookieSecure != test.expected {
				t.Errorf("%s: Expected CookieSecure %v, got %v", test.description, test.expected, config.CookieSecure)
			}
		})
	}
}

func TestApplyEnvConfiguration_Combined(t *testing.T) {
	os.Setenv("J_SESSION_KEY", "abcdefghijklmnopqrstuvwxyz123456")
	os.Setenv("J_SESSION_NAME", "my-app-session")
	os.Setenv("J_COOKIE_DOMAIN", ".myapp.com")
	os.Setenv("J_COOKIE_MAX_AGE", "1800")
	os.Setenv("J_COOKIE_HTTPONLY", "0")
	os.Setenv("J_SSL_CERT", "/path/to/cert.pem")
	os.Setenv("J_PORT", "8080")
	defer func() {
		os.Unsetenv("J_SESSION_KEY")
		os.Unsetenv("J_SESSION_NAME")
		os.Unsetenv("J_COOKIE_DOMAIN")
		os.Unsetenv("J_COOKIE_MAX_AGE")
		os.Unsetenv("J_COOKIE_HTTPONLY")
		os.Unsetenv("J_SSL_CERT")
		os.Unsetenv("J_PORT")
	}()

	config := DefaultConfiguration()
	ApplyEnvConfiguration(&config)

	if config.SessionKey != "abcdefghijklmnopqrstuvwxyz123456" {
		t.Errorf("Expected SessionKey 'abcdefghijklmnopqrstuvwxyz123456', got %q", config.SessionKey)
	}
	if config.SessionName != "my-app-session" {
		t.Errorf("Expected SessionName 'my-app-session', got %q", config.SessionName)
	}
	if config.CookieDomain != ".myapp.com" {
		t.Errorf("Expected CookieDomain '.myapp.com', got %q", config.CookieDomain)
	}
	if config.CookieMaxAge != 1800 {
		t.Errorf("Expected CookieMaxAge 1800, got %d", config.CookieMaxAge)
	}
	if config.CookieHTTPOnly != false {
		t.Errorf("Expected CookieHTTPOnly false, got %v", config.CookieHTTPOnly)
	}
	if config.CookieSecure != true {
		t.Errorf("Expected CookieSecure true (SSL enabled), got %v", config.CookieSecure)
	}
	if config.Port != "8080" {
		t.Errorf("Expected Port '8080', got %q", config.Port)
	}
}
