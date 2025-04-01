package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateTempDir creates a temporary directory for testing and returns its path
func CreateTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "awsd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// CreateMockAWSConfig creates a mock AWS credentials file with test profiles
func CreateMockAWSConfig(t *testing.T, dir string) string {
	t.Helper()
	config := `[default]
aws_access_key_id = test_access_key
aws_secret_access_key = test_secret_key
region = us-east-1

[profile dev]
aws_access_key_id = dev_access_key
aws_secret_access_key = dev_secret_key
region = us-west-2

[profile prod]
aws_access_key_id = prod_access_key
aws_secret_access_key = prod_secret_key
region = eu-west-1`

	configPath := filepath.Join(dir, "config")
	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		t.Fatalf("Failed to write mock AWS config: %v", err)
	}
	return configPath
}

// CleanupTempDir removes a temporary directory and its contents
func CleanupTempDir(t *testing.T, dir string) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup temp dir: %v", err)
	}
}

// SetTestEnv sets up test environment variables
func SetTestEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set environment variable %s: %v", key, err)
	}
}

// UnsetTestEnv removes test environment variables
func UnsetTestEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Errorf("Failed to unset environment variable %s: %v", key, err)
	}
}
