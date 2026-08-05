package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// version is injected at build time via ldflags:
//
//	-X github.com/radiusmethod/awsd/src/cmd.version=v1.2.3
//
// GoReleaser sets it from the git tag, and `make install` sets it from
// `git describe`. Plain `go build` leaves it as "dev".
var version = "dev"

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
