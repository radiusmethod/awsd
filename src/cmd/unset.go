package cmd

import (
	"log"

	"github.com/radiusmethod/awsd/src/utils"
	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset the active AWS profile or region.",
	Long:  "Clear the active AWS profile (back to default) or AWS region.",
}

var unsetProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Reset the active AWS profile to default.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := utils.GetHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		if err := utils.WriteFile("default", homeDir); err != nil {
			log.Fatal(err)
		}
		printColoredMessage("Profile reset to default.\n", utils.PromptColor)
	},
}

var unsetRegionCmd = &cobra.Command{
	Use:   "region",
	Short: "Clear the active AWS region.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := utils.GetHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		if err := utils.UnsetRegion(homeDir); err != nil {
			log.Fatal(err)
		}
		printColoredMessage("Region cleared.\n", utils.PromptColor)
	},
}

func init() {
	unsetCmd.AddCommand(unsetProfileCmd)
	unsetCmd.AddCommand(unsetRegionCmd)
	rootCmd.AddCommand(unsetCmd)
}
