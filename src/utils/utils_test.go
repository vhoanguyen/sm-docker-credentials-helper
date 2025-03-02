package utils

import (
	"os"
	"testing"
)

func TestParseEnv(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("TEST_ENV_1", "value1")
	os.Setenv("TEST_ENV_2", "value2")
	defer os.Unsetenv("TEST_ENV_1")
	defer os.Unsetenv("TEST_ENV_2")

	tests := []struct {
		name     string
		envs     []string
		expected []string
	}{
		{
			name:     "All environment variables set",
			envs:     []string{"TEST_ENV_1", "TEST_ENV_2"},
			expected: []string{"value1", "value2"},
		},
		{
			name:     "One environment variable not set",
			envs:     []string{"TEST_ENV_1", "TEST_ENV_3"},
			expected: []string{"value1", ""},
		},
		{
			name:     "No environment variables set",
			envs:     []string{"TEST_ENV_3", "TEST_ENV_4"},
			expected: []string{"", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseEnv(tt.envs)
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected %s, got %s", tt.expected[i], v)
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "Item in slice",
			slice:    []string{"a", "b", "c"},
			item:     "b",
			expected: true,
		},
		{
			name:     "Item not in slice",
			slice:    []string{"a", "b", "c"},
			item:     "d",
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []string{},
			item:     "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}