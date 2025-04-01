package cmd

import (
	"testing"

	"github.com/radiusmethod/awsd/src/utils/testutils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	cmd := listCmd
	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "List AWS profiles command.", cmd.Short)
	assert.Equal(t, "This lists all your AWS profiles.", cmd.Long)
	assert.Equal(t, []string{"l"}, cmd.Aliases)
}

func TestRunProfileLister(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	// Create mock AWS config
	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	testutils.SetTestEnv(t, "AWS_CONFIG_FILE", configPath)
	defer testutils.UnsetTestEnv(t, "AWS_CONFIG_FILE")

	err := runProfileLister()
	assert.NoError(t, err)
}

func TestListCommandIntegration(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	// Create mock AWS config
	configPath := testutils.CreateMockAWSConfig(t, tempDir)
	testutils.SetTestEnv(t, "AWS_CONFIG_FILE", configPath)
	defer testutils.UnsetTestEnv(t, "AWS_CONFIG_FILE")

	// Set HOME environment variable
	testutils.SetTestEnv(t, "HOME", tempDir)
	defer testutils.UnsetTestEnv(t, "HOME")

	// Create a new command instance
	cmd := &cobra.Command{
		Use:   "awsd",
		Short: "awsd - switch between AWS profiles.",
		Long:  "Allows for switching AWS profiles files.",
	}

	// Add the list command
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List AWS profiles command.",
		Aliases: []string{"l"},
		Long:    "This lists all your AWS profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			err := runProfileLister()
			if err != nil {
				t.Fatal(err)
			}
		},
	}
	cmd.AddCommand(listCmd)

	// Test both list and l aliases
	aliases := []string{"list", "l"}
	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			cmd.SetArgs([]string{alias})
			err := cmd.Execute()
			assert.NoError(t, err)
		})
	}
}
