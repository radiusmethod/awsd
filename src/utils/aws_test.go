package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radiusmethod/awsd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment(t *testing.T) (string, func()) {
	// Save original environment
	origHome := os.Getenv("HOME")
	origConfigFile := os.Getenv("AWS_CONFIG_FILE")

	// Create temporary directory
	tempDir := testutils.CreateTempDir(t)

	// Create .aws directory
	awsDir := filepath.Join(tempDir, ".aws")
	err := os.MkdirAll(awsDir, 0755)
	assert.NoError(t, err)

	// Set environment variables
	os.Setenv("HOME", tempDir)
	os.Setenv("AWS_CONFIG_FILE", filepath.Join(tempDir, "config"))

	// Return cleanup function
	cleanup := func() {
		os.Setenv("HOME", origHome)
		os.Setenv("AWS_CONFIG_FILE", origConfigFile)
		testutils.CleanupTempDir(t, tempDir)
	}

	return tempDir, cleanup
}

func TestGetProfiles(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a mock AWS config file
	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	os.Setenv("AWS_CONFIG_FILE", configPath)

	// Test getting profiles
	profiles, err := GetProfiles()
	assert.NoError(t, err, "Should not return error")

	// Verify the profiles
	expectedProfiles := []string{"default", "dev", "prod"}
	assert.Equal(t, expectedProfiles, profiles, "Expected profiles should match")
}

func TestGetProfilesWithComplexConfig(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Copy complex config to temp directory
	complexConfigPath := filepath.Join("..", "..", "testdata", "aws_config_examples", "complex_config")
	configContent, err := os.ReadFile(complexConfigPath)
	if err != nil {
		t.Fatalf("Failed to read complex config: %v", err)
	}

	configPath := filepath.Join(tempDir, "config")
	if err := os.WriteFile(configPath, configContent, 0600); err != nil {
		t.Fatalf("Failed to write complex config: %v", err)
	}

	// Test getting profiles
	profiles, err := GetProfiles()
	assert.NoError(t, err, "Should not return error")

	// Verify the profiles
	expectedProfiles := []string{
		"default",
		"dev",
		"dev.admin",
		"dev.readonly",
		"prod",
		"prod.admin",
		"prod.readonly",
	}
	assert.Equal(t, expectedProfiles, profiles, "Expected profiles should match")
}

func TestGetProfilesWithMalformedConfig(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a malformed config file
	configPath := filepath.Join(tempDir, "config")
	malformedContent := `[default]
aws_access_key_id = test_access_key
aws_secret_access_key = test_secret_key
region = us-east-1

[profile dev
aws_access_key_id = dev_access_key
aws_secret_access_key = dev_secret_key
region = us-west-2`

	err := os.WriteFile(configPath, []byte(malformedContent), 0600)
	assert.NoError(t, err)

	// Test getting profiles with error handling
	profiles, err := GetProfiles()
	assert.Error(t, err, "Should return an error for malformed config")
	assert.Nil(t, profiles, "Should return nil profiles for malformed config")
	assert.Contains(t, err.Error(), "unclosed section", "Error should mention unclosed section")
}

func TestGetProfilesWithError(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test with valid config
	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	os.Setenv("AWS_CONFIG_FILE", configPath)

	profiles, err := GetProfiles()
	assert.NoError(t, err, "Should not return error for valid config")
	expectedProfiles := []string{"default", "dev", "prod"}
	assert.Equal(t, expectedProfiles, profiles, "Expected profiles should match")

	// Test with malformed config
	malformedContent := `[default]
aws_access_key_id = test_access_key
aws_secret_access_key = test_secret_key
region = us-east-1

[profile dev
aws_access_key_id = dev_access_key
aws_secret_access_key = dev_secret_key
region = us-west-2`

	err = os.WriteFile(configPath, []byte(malformedContent), 0600)
	assert.NoError(t, err)

	profiles, err = GetProfiles()
	assert.Error(t, err, "Should return error for malformed config")
	assert.Nil(t, profiles, "Should return nil profiles for malformed config")
	assert.Contains(t, err.Error(), "unclosed section", "Error should mention unclosed section")

	// Test with non-existent config file
	os.Setenv("AWS_CONFIG_FILE", "/nonexistent/config")
	profiles, err = GetProfiles()
	assert.Error(t, err, "Should return error for non-existent config")
	assert.Nil(t, profiles, "Should return nil profiles for non-existent config")
}
