package env

import (
    "os"
    "path/filepath"
    "testing"
)

func TestParse(t *testing.T) {
    tests := []struct {
        name     string
        content  string
        expected map[string]string
    }{
        {
            name: "basic key-value pairs",
            content: `KEY1=value1
KEY2=value2
KEY3=value3`,
            expected: map[string]string{
                "KEY1": "value1",
                "KEY2": "value2",
                "KEY3": "value3",
            },
        },
        {
            name: "with comments",
            content: `# This is a comment
KEY1=value1
# Another comment
KEY2=value2`,
            expected: map[string]string{
                "KEY1": "value1",
                "KEY2": "value2",
            },
        },
        {
            name: "with empty lines",
            content: `KEY1=value1

KEY2=value2

`,
            expected: map[string]string{
                "KEY1": "value1",
                "KEY2": "value2",
            },
        },
        {
            name: "with quoted values",
            content: `KEY1="value with spaces"
KEY2='single quoted value'
KEY3=unquoted`,
            expected: map[string]string{
                "KEY1": "value with spaces",
                "KEY2": "single quoted value",
                "KEY3": "unquoted",
            },
        },
        {
            name: "with spaces around equals",
            content: `KEY1 = value1
KEY2= value2
KEY3 =value3`,
            expected: map[string]string{
                "KEY1": "value1",
                "KEY2": "value2",
                "KEY3": "value3",
            },
        },
        {
            name: "with equals in value",
            content: `KEY1=value=with=equals
KEY2=http://example.com?param=value`,
            expected: map[string]string{
                "KEY1": "value=with=equals",
                "KEY2": "http://example.com?param=value",
            },
        },
        {
            name: "malformed lines are skipped",
            content: `KEY1=value1
INVALID_LINE_NO_EQUALS
KEY2=value2`,
            expected: map[string]string{
                "KEY1": "value1",
                "KEY2": "value2",
            },
        },
        {
            name:     "empty file",
            content:  "",
            expected: map[string]string{},
        },
        {
            name: "only comments and empty lines",
            content: `# Comment 1
# Comment 2

# Comment 3`,
            expected: map[string]string{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a temporary .env file
            tmpDir := t.TempDir()
            envFile := filepath.Join(tmpDir, ".env")

            if err := os.WriteFile(envFile, []byte(tt.content), 0644); err != nil {
                t.Fatalf("Failed to create temp file: %v", err)
            }

            // Parse the file
            result, err := Parse(envFile)
            if err != nil {
                t.Fatalf("Parse() error = %v", err)
            }

            // Check the results
            if len(result) != len(tt.expected) {
                t.Errorf("Expected %d entries, got %d", len(tt.expected), len(result))
            }

            for key, expectedValue := range tt.expected {
                if actualValue, ok := result[key]; !ok {
                    t.Errorf("Missing key %q", key)
                } else if actualValue != expectedValue {
                    t.Errorf("For key %q: expected %q, got %q", key, expectedValue, actualValue)
                }
            }

            for key := range result {
                if _, ok := tt.expected[key]; !ok {
                    t.Errorf("Unexpected key %q with value %q", key, result[key])
                }
            }
        })
    }
}

func TestParseNonExistentFile(t *testing.T) {
    // Parsing a non-existent file should return an empty map, not an error
    result, err := Parse("/nonexistent/path/.env")
    if err != nil {
        t.Errorf("Parse() should not error on non-existent file, got: %v", err)
    }
    if len(result) != 0 {
        t.Errorf("Expected empty map, got %d entries", len(result))
    }
}

func TestParseInvalidPath(t *testing.T) {
    // A path with a null byte is invalid on all platforms, returning non-IsNotExist error
    _, err := Parse("/tmp/test\x00file")
    if err == nil {
        t.Error("Expected error for path with embedded null byte")
    }
}

func TestUnquote(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {`"double quoted"`, "double quoted"},
        {`'single quoted'`, "single quoted"},
        {`unquoted`, "unquoted"},
        {`"`, `"`},
        {`''`, ``},
        {`""`, ``},
        {`"mismatched'`, `"mismatched'`},
        {`'mismatched"`, `'mismatched"`},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result := unquote(tt.input)
            if result != tt.expected {
                t.Errorf("unquote(%q) = %q, expected %q", tt.input, result, tt.expected)
            }
        })
    }
}
