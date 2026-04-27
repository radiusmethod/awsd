package cmd

import (
	"fmt"
	"log"

	"github.com/radiusmethod/awsd/src/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:       "list [profiles|regions]",
	Short:     "List AWS profiles or regions.",
	Aliases:   []string{"l"},
	Long:      "List all your AWS profiles, or the known AWS regions. Defaults to profiles.",
	ValidArgs: []string{"profiles", "regions"},
	Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		target := "profiles"
		if len(args) == 1 {
			target = args[0]
		}
		var err error
		switch target {
		case "regions":
			err = runRegionLister()
		default:
			err = runProfileLister()
		}
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runProfileLister() error {
	profiles, err := utils.GetProfiles()
	if err != nil {
		return err
	}
	for _, p := range profiles {
		fmt.Println(p)
	}
	return nil
}

func runRegionLister() error {
	for _, r := range utils.GetRegions() {
		fmt.Println(r)
	}
	return nil
}
