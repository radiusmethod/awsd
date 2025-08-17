package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/radiusmethod/awsd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTouchFile(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	filePath := filepath.Join(tempDir, "test.txt")
	err := TouchFile(filePath)
	assert.NoError(t, err, "Should create file without error")

	// Verify file exists
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "File should exist")
}

func TestWriteFile(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	// Test writing a profile
	err := WriteFile("test-profile", tempDir)
	assert.NoError(t, err, "Should write file without error")

	// Verify file exists and contains correct content
	filePath := filepath.Join(tempDir, ".awsd")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Should read file without error")
	assert.Equal(t, "test-profile", string(content), "File content should match")

	// Test writing default profile
	err = WriteFile("default", tempDir)
	assert.NoError(t, err, "Should write default profile without error")

	content, err = os.ReadFile(filePath)
	assert.NoError(t, err, "Should read file without error")
	assert.Equal(t, "", string(content), "Default profile should write empty string")
}

func TestGetEnv(t *testing.T) {
	// Test with environment variable set
	testutils.SetTestEnv(t, "TEST_VAR", "test-value")
	defer testutils.UnsetTestEnv(t, "TEST_VAR")

	value := GetEnv("TEST_VAR", "fallback")
	assert.Equal(t, "test-value", value, "Should return environment variable value")

	// Test with environment variable not set
	value = GetEnv("NONEXISTENT_VAR", "fallback")
	assert.Equal(t, "fallback", value, "Should return fallback value")
}

func TestGetHomeDir(t *testing.T) {
	// Save original HOME value
	origHome := os.Getenv("HOME")
	defer func() {
		if origHome != "" {
			os.Setenv("HOME", origHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// Test with HOME set
	testutils.SetTestEnv(t, "HOME", "/test/home")
	homeDir, err := GetHomeDir()
	assert.NoError(t, err, "Should not return error when HOME is set")
	assert.Equal(t, "/test/home", homeDir, "Should return HOME environment variable value")

	// Test with HOME unset
	testutils.UnsetTestEnv(t, "HOME")
	homeDir, err = GetHomeDir()
	if err != nil {
		// This is okay - on some systems UserHomeDir() might fail when HOME is unset
		assert.Contains(t, err.Error(), "error getting user home directory")
	} else {
		assert.NotEmpty(t, homeDir, "Should return UserHomeDir value")
	}
}

func TestGetProfileFileLocation(t *testing.T) {
	// Save original HOME value
	origHome := os.Getenv("HOME")
	defer func() {
		if origHome != "" {
			os.Setenv("HOME", origHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	// Set HOME environment variable
	testutils.SetTestEnv(t, "HOME", tempDir)
	expectedPath := filepath.Join(tempDir, ".aws")
	actualPath := GetProfileFileLocation()
	assert.Equal(t, expectedPath, actualPath, "Should return correct .aws directory path")
}

func TestGetCurrentProfileFile(t *testing.T) {
	// Test with AWS_CONFIG_FILE set
	testutils.SetTestEnv(t, "AWS_CONFIG_FILE", "/test/config")
	defer testutils.UnsetTestEnv(t, "AWS_CONFIG_FILE")

	file := GetCurrentProfileFile()
	assert.Equal(t, "/test/config", file, "Should return AWS_CONFIG_FILE value")

	// Test without AWS_CONFIG_FILE set
	testutils.UnsetTestEnv(t, "AWS_CONFIG_FILE")
	file = GetCurrentProfileFile()
	assert.Contains(t, file, ".aws/config", "Should return default config path")
}

func TestIsDirectoryExists(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	// Test existing directory
	exists := IsDirectoryExists(tempDir)
	assert.True(t, exists, "Should return true for existing directory")

	// Test non-existing directory
	exists = IsDirectoryExists("/nonexistent/directory")
	assert.False(t, exists, "Should return false for non-existing directory")
}

func TestAppendIfNotExists(t *testing.T) {
	// Test appending new item
	slice := []string{"a", "b", "c"}
	result := AppendIfNotExists(slice, "d")
	assert.Equal(t, []string{"a", "b", "c", "d"}, result, "Should append new item")

	// Test appending existing item
	result = AppendIfNotExists(slice, "b")
	assert.Equal(t, []string{"a", "b", "c"}, result, "Should not append existing item")

	// Test appending to empty slice
	var emptySlice []string
	result = AppendIfNotExists(emptySlice, "a")
	assert.Equal(t, []string{"a"}, result, "Should append to empty slice")
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	// Test existing item
	assert.True(t, Contains(slice, "b"), "Should find existing item")

	// Test non-existing item
	assert.False(t, Contains(slice, "d"), "Should not find non-existing item")

	// Test empty slice
	var emptySlice []string
	assert.False(t, Contains(emptySlice, "a"), "Should return false for empty slice")
}

func TestCheckError(t *testing.T) {
	// Test ^D error
	if os.Getenv("TEST_CHECK_ERROR_DEL") == "1" {
		CheckError(fmt.Errorf("^D"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckError")
	cmd.Env = append(os.Environ(), "TEST_CHECK_ERROR_DEL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}

func TestCheckErrorCtrlC(t *testing.T) {
	// Test ^C error
	if os.Getenv("TEST_CHECK_ERROR_CTRL_C") == "1" {
		CheckError(fmt.Errorf("^C"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErrorCtrlC")
	cmd.Env = append(os.Environ(), "TEST_CHECK_ERROR_CTRL_C=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}

func TestCheckErrorOther(t *testing.T) {
	// Test other error
	if os.Getenv("TEST_CHECK_ERROR_OTHER") == "1" {
		CheckError(fmt.Errorf("other error"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErrorOther")
	cmd.Env = append(os.Environ(), "TEST_CHECK_ERROR_OTHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}

func TestBellSkipper(t *testing.T) {
	bs := &BellSkipper{}

	// Test writing bell character
	n, err := bs.Write([]byte{7}) // ASCII bell character
	assert.NoError(t, err, "Write should not return error")
	assert.Equal(t, 0, n, "Write should skip bell character")

	// Test writing normal text
	n, err = bs.Write([]byte("test"))
	assert.NoError(t, err, "Write should not return error")
	assert.Equal(t, 4, n, "Write should return correct number of bytes written")
}

func TestNewPromptUISearcher(t *testing.T) {
	items := []string{"test1", "test2", "another", "something"}
	searcher := NewPromptUISearcher(items)

	// Test exact match
	result := searcher("test1", 0)
	assert.True(t, result, "Should find exact match")

	// Test partial match
	result = searcher("test", 0)
	assert.True(t, result, "Should find partial match")

	// Test case insensitive match
	result = searcher("TEST1", 0)
	assert.True(t, result, "Should find case insensitive match")

	// Test no match
	result = searcher("nonexistent", 0)
	assert.False(t, result, "Should not find non-existent item")

	// Test empty search
	result = searcher("", 0)
	assert.True(t, result, "Should match on empty search")

	// Test different index
	result = searcher("test", 1)
	assert.True(t, result, "Should find match at different index")

	// Test index at boundary
	result = searcher("test", len(items)-1)
	assert.False(t, result, "Should handle index at boundary")
}
