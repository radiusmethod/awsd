package cmd

import (
	"fmt"
	"log"

	"github.com/radiusmethod/awsd/src/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List AWS profiles command.",
	Aliases: []string{"l"},
	Long:    "This lists all your AWS profiles.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runProfileLister()
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
