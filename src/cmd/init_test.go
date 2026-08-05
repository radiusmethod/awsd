package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/radiusmethod/awsd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

// captureStream runs fn with the given stream redirected and returns what was
// written to it. The commands print with fmt.Print rather than through cobra's
// writer, since their output is meant to be eval'd by the shell.
func captureStream(t *testing.T, stream **os.File, fn func()) string {
	t.Helper()
	orig := *stream
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	*stream = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe: %v", err)
	}
	*stream = orig

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to read pipe: %v", err)
	}
	return buf.String()
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	return captureStream(t, &os.Stdout, fn)
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	return captureStream(t, &os.Stderr, fn)
}

func TestShellenvCommand(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)
	t.Setenv("HOME", tempDir)

	tests := []struct {
		name     string
		awsdFile string
		args     []string
		expected string
	}{
		{
			name:     "Profile and region, default shell",
			awsdFile: "profile=dev\nregion=us-east-1\n",
			args:     []string{},
			expected: "export AWS_PROFILE='dev'\nexport AWS_REGION='us-east-1'\nexport AWS_DEFAULT_REGION='us-east-1'\n",
		},
		{
			name:     "Explicit shell argument",
			awsdFile: "profile=dev\n",
			args:     []string{"fish"},
			expected: "set -gx AWS_PROFILE 'dev'\n",
		},
		{
			name:     "Empty profile unsets",
			awsdFile: "profile=\n",
			args:     []string{"bash"},
			expected: "unset AWS_PROFILE\n",
		},
		{
			name:     "Legacy single-line file",
			awsdFile: "dev\n",
			args:     []string{"bash"},
			expected: "export AWS_PROFILE='dev'\n",
		},
		{
			name:     "No region line leaves region alone",
			awsdFile: "profile=dev\n",
			args:     []string{"bash"},
			expected: "export AWS_PROFILE='dev'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			awsdPath := filepath.Join(tempDir, ".awsd")
			if err := os.WriteFile(awsdPath, []byte(tt.awsdFile), 0644); err != nil {
				t.Fatalf("Failed to write .awsd: %v", err)
			}
			out := captureStdout(t, func() {
				shellenvCmd.Run(shellenvCmd, tt.args)
			})
			assert.Equal(t, tt.expected, out)
		})
	}
}

func TestShellenvCommandNoStateFile(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)
	t.Setenv("HOME", tempDir)

	out := captureStdout(t, func() {
		shellenvCmd.Run(shellenvCmd, []string{"bash"})
	})
	assert.Equal(t, "unset AWS_PROFILE\n", out)
}

func TestInitCommand(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", "powershell", "pwsh"} {
		t.Run(shell, func(t *testing.T) {
			out := captureStdout(t, func() {
				initCmd.Run(initCmd, []string{shell})
			})
			assert.NotEmpty(t, out)
			assert.Contains(t, out, "awsd")
		})
	}
}

func TestInitCommandArgValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{name: "Valid shell", args: []string{"zsh"}},
		{name: "Unsupported shell", args: []string{"csh"}, expectError: true},
		{name: "No shell", args: []string{}, expectError: true},
		{name: "Too many args", args: []string{"zsh", "bash"}, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := initCmd.Args(initCmd, tt.args)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
