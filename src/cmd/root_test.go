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
			name:     "Set command",
			args:     []string{"awsd", "set"},
			expected: false,
		},
		{
			name:     "Unset command",
			args:     []string{"awsd", "unset"},
			expected: false,
		},
		{
			name:     "Init command",
			args:     []string{"awsd", "init"},
			expected: false,
		},
		{
			name:     "Shellenv command",
			args:     []string{"awsd", "shellenv"},
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

	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	t.Setenv("AWS_CONFIG_FILE", configPath)
	t.Setenv("HOME", tempDir)

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
			expectContent: "profile=dev\n",
		},
		{
			name:          "Invalid profile",
			profile:       "invalid",
			expectError:   true,
			expectFile:    false,
			expectContent: "",
		},
		{
			name:          "Default profile",
			profile:       "default",
			expectError:   false,
			expectFile:    true,
			expectContent: "profile=\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			awsdFile := filepath.Join(tempDir, ".awsd")
			_ = os.Remove(awsdFile)

			err := directProfileSwitch(tt.profile)
			if tt.expectError {
				assert.ErrorIs(t, err, errProfileNotFound)
				_, statErr := os.Stat(awsdFile)
				assert.True(t, os.IsNotExist(statErr), "File should not exist for invalid profile")
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

// The shell integration evals command substitutions of this binary, so a
// warning on stdout gets executed rather than shown. An older binary printing
// "Profile init does not exist" to stdout is what turns a plain version skew
// into "(eval):1: bad pattern: ^[[0".
func TestDirectProfileSwitchWarningGoesToStderr(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	t.Setenv("AWS_CONFIG_FILE", configPath)
	t.Setenv("HOME", tempDir)

	var stderr string
	stdout := captureStdout(t, func() {
		stderr = captureStderr(t, func() {
			assert.ErrorIs(t, directProfileSwitch("invalid"), errProfileNotFound)
		})
	})

	assert.Empty(t, stdout, "warning must not reach stdout, the shell integration evals it")
	assert.Contains(t, stderr, "does not exist")
}

// A successful switch still reports on stdout.
func TestDirectProfileSwitchSuccessGoesToStdout(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	t.Setenv("AWS_CONFIG_FILE", configPath)
	t.Setenv("HOME", tempDir)

	stdout := captureStdout(t, func() {
		assert.NoError(t, directProfileSwitch("dev"))
	})
	assert.Contains(t, stdout, "dev")
}

func TestRootCommand(t *testing.T) {
	cmd := rootCmd
	assert.NotNil(t, cmd)
	assert.Equal(t, "awsd", cmd.Use)
	assert.Equal(t, "awsd - switch between AWS profiles.", cmd.Short)
	assert.Equal(t, "Allows for switching AWS profiles files.", cmd.Long)
}

func TestPrintColoredMessage(t *testing.T) {
	printColoredMessage("test", "test")
}
