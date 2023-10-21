package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version string = "v0.0.7"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "awsd version command",
	Aliases: []string{"v"},
	Long:    "Returns the current version of awsd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("awsd version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
