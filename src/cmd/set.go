package cmd

import (
	"fmt"
	"log"

	"github.com/radiusmethod/awsd/src/utils"
	"github.com/radiusmethod/promptui"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the active AWS profile or region.",
	Long:  "Set the active AWS profile or region. With no subcommand, prints help.",
}

var setProfileCmd = &cobra.Command{
	Use:   "profile [name]",
	Short: "Set the active AWS profile (interactive picker if no name given).",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := runProfileSwitcher(); err != nil {
				log.Fatal(err)
			}
			return
		}
		handleSwitchError(directProfileSwitch(args[0]))
	},
}

var setRegionCmd = &cobra.Command{
	Use:   "region [region]",
	Short: "Set the active AWS region (interactive picker if no region given).",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var region string
		if len(args) == 0 {
			r, err := runRegionSwitcher()
			if err != nil {
				log.Fatal(err)
			}
			if r == "" {
				return
			}
			region = r
		} else {
			region = args[0]
		}
		homeDir, err := utils.GetHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		if err := utils.WriteRegion(region, homeDir); err != nil {
			log.Fatal(err)
		}
		printColoredMessage("Region ", utils.PromptColor)
		printColoredMessage(region, utils.CyanColor)
		printColoredMessage(" set.\n", utils.PromptColor)
	},
}

func init() {
	setCmd.AddCommand(setProfileCmd)
	setCmd.AddCommand(setRegionCmd)
	rootCmd.AddCommand(setCmd)
}

func runRegionSwitcher() (string, error) {
	regions := utils.GetRegions()
	fmt.Printf(utils.NoticeColor, "AWS Region Switcher\n")
	region, err := getRegionFromPrompt(regions)
	if err != nil {
		return "", err
	}
	fmt.Printf(utils.PromptColor, "Choose a region")
	fmt.Printf(utils.NoticeColor, "? ")
	fmt.Printf(utils.CyanColor, region)
	fmt.Println()
	return region, nil
}

func getRegionFromPrompt(regions []string) (string, error) {
	prompt := promptui.Select{
		Label:        fmt.Sprintf(utils.PromptColor, "Choose a region"),
		Items:        regions,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | cyan }}",
		},
		Searcher:          utils.NewPromptUISearcher(regions),
		StartInSearchMode: true,
		Stdout:            &utils.BellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.CheckError(err)
		return "", nil
	}
	return result, nil
}
