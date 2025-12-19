package client

import (
	"testing"
)

func TestSubstituteEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "No variables",
			input:    "http://example.com/path",
			envVars:  map[string]string{},
			expected: "http://example.com/path",
		},
		{
			name:  "Single variable with braces",
			input: "http://${HOST}/path",
			envVars: map[string]string{
				"HOST": "example.com",
			},
			expected: "http://example.com/path",
		},
		{
			name:  "Single variable without braces",
			input: "http://$HOST/path",
			envVars: map[string]string{
				"HOST": "example.com",
			},
			expected: "http://example.com/path",
		},
		{
			name:  "Multiple variables",
			input: "${PROTO}://${HOST}:${PORT}/path",
			envVars: map[string]string{
				"PROTO": "https",
				"HOST":  "example.com",
				"PORT":  "8080",
			},
			expected: "https://example.com:8080/path",
		},
		{
			name:     "Variable not set",
			input:    "http://${NOTSET}/path",
			envVars:  map[string]string{},
			expected: "http://${NOTSET}/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			result := substituteEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("substituteEnvVars() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient() returned nil")
	}
	if client.httpClient == nil {
		t.Error("NewClient() returned client with nil httpClient")
	}
}
