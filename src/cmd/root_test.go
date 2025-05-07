package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radiusmethod/awsd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestShouldRunDirectProfileSwitch(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "Direct profile switch",
			args:     []string{"awsd", "dev"},
			expected: true,
		},
		{
			name:     "List command",
			args:     []string{"awsd", "list"},
			expected: false,
		},
		{
			name:     "Help command",
			args:     []string{"awsd", "--help"},
			expected: false,
		},
		{
			name:     "Version command",
			args:     []string{"awsd", "version"},
			expected: false,
		},
		{
			name:     "No arguments",
			args:     []string{"awsd"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			result := shouldRunDirectProfileSwitch()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDirectProfileSwitch(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	// Create mock AWS config
	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	testutils.SetTestEnv(t, "AWS_CONFIG_FILE", configPath)
	defer testutils.UnsetTestEnv(t, "AWS_CONFIG_FILE")

	// Set HOME environment variable
	testutils.SetTestEnv(t, "HOME", tempDir)
	defer testutils.UnsetTestEnv(t, "HOME")

	tests := []struct {
		name          string
		profile       string
		expectError   bool
		expectFile    bool
		expectContent string
	}{
		{
			name:          "Valid profile",
			profile:       "dev",
			expectError:   false,
			expectFile:    true,
			expectContent: "dev",
		},
		{
			name:          "Invalid profile",
			profile:       "invalid",
			expectError:   false,
			expectFile:    false,
			expectContent: "",
		},
		{
			name:          "Default profile",
			profile:       "default",
			expectError:   false,
			expectFile:    true,
			expectContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Remove any existing .awsd file before each test
			awsdFile := filepath.Join(tempDir, ".awsd")
			_ = os.Remove(awsdFile)

			err := directProfileSwitch(tt.profile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tt.expectFile {
				content, err := os.ReadFile(awsdFile)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectContent, string(content))
			} else {
				_, err := os.Stat(awsdFile)
				assert.True(t, os.IsNotExist(err), "File should not exist for invalid profile")
			}
		})
	}
}

func TestRootCommand(t *testing.T) {
	cmd := rootCmd
	assert.NotNil(t, cmd)
	assert.Equal(t, "awsd", cmd.Use)
	assert.Equal(t, "awsd - switch between AWS profiles.", cmd.Short)
	assert.Equal(t, "Allows for switching AWS profiles files.", cmd.Long)
}

func TestPrintColoredMessage(t *testing.T) {
	// This is a simple test that just ensures the function doesn't panic
	// since we can't easily capture stdout in tests
	printColoredMessage("test", "test")
}
