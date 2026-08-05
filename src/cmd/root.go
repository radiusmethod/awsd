package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/radiusmethod/awsd/src/utils"
	"github.com/radiusmethod/promptui"
	"github.com/spf13/cobra"
)

// errProfileNotFound is returned when argv named a profile that isn't in the
// AWS config. The generated shell function keys off the exit code, so this has
// to fail rather than exit 0.
var errProfileNotFound = errors.New("profile does not exist")

var rootCmd = &cobra.Command{
	Use:   "awsd",
	Short: "awsd - switch between AWS profiles.",
	Long:  "Allows for switching AWS profiles files.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runProfileSwitcher(); err != nil {
			log.Fatal(err)
		}
	},
}

// RootCmd returns the root cobra command. Used by docs generation.
func RootCmd() *cobra.Command {
	return rootCmd
}

// Entry point for the CLI tool
func Execute() {
	if shouldRunDirectProfileSwitch() {
		profile := os.Args[1]
		handleSwitchError(directProfileSwitch(profile))
		return
	}
	runRootCmd()
}

func runRootCmd() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runProfileSwitcher() error {
	profiles, err := utils.GetProfiles()
	if err != nil {
		return err
	}
	fmt.Printf(utils.NoticeColor, "AWS Profile Switcher\n")
	profile, err := getProfileFromPrompt(profiles)
	if err != nil {
		return err
	}
	fmt.Printf(utils.PromptColor, "Choose a profile")
	fmt.Printf(utils.NoticeColor, "? ")
	fmt.Printf(utils.CyanColor, profile)
	fmt.Println()
	homeDir, err := utils.GetHomeDir()
	if err != nil {
		return err
	}
	return utils.WriteFile(profile, homeDir)
}

func shouldRunDirectProfileSwitch() bool {
	// Any argv[1] not in this list is treated as a profile name, so every
	// subcommand and alias has to be listed here. Escape hatch for a profile
	// that collides with one of these: awsd set profile <name>.
	invalidProfiles := []string{"l", "list", "set", "unset", "init", "shellenv", "completion", "help", "--help", "-h", "v", "version"}
	return len(os.Args) > 1 && !utils.Contains(invalidProfiles, os.Args[1])
}

func directProfileSwitch(desiredProfile string) error {
	profiles, err := utils.GetProfiles()
	if err != nil {
		return err
	}
	if utils.Contains(profiles, desiredProfile) {
		printColoredMessage("Profile ", utils.PromptColor)
		printColoredMessage(desiredProfile, utils.CyanColor)
		printColoredMessage(" set.\n", utils.PromptColor)
		homeDir, err := utils.GetHomeDir()
		if err != nil {
			return err
		}
		return utils.WriteFile(desiredProfile, homeDir)
	}
	printColoredWarning("WARNING: Profile ", utils.NoticeColor)
	printColoredWarning(desiredProfile, utils.CyanColor)
	printColoredWarning(" does not exist.\n", utils.PromptColor)
	return fmt.Errorf("%w: %s", errProfileNotFound, desiredProfile)
}

// handleSwitchError exits non-zero when a profile switch fails. A missing
// profile has already been reported on stderr, so exit quietly instead of
// printing it a second time through log.Fatal.
func handleSwitchError(err error) {
	if err == nil {
		return
	}
	if errors.Is(err, errProfileNotFound) {
		os.Exit(1)
	}
	log.Fatal(err)
}

func getProfileFromPrompt(profiles []string) (string, error) {
	prompt := promptui.Select{
		Label:        fmt.Sprintf(utils.PromptColor, "Choose a profile"),
		Items:        profiles,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | cyan }}",
		},
		Searcher:          utils.NewPromptUISearcher(profiles),
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

func printColoredMessage(msg, color string) {
	fmt.Printf(color, msg)
}

// printColoredWarning writes to stderr. Anything the shell integration might
// capture in a command substitution has to stay off stdout, or the caller ends
// up trying to eval ANSI escapes.
func printColoredWarning(msg, color string) {
	fmt.Fprintf(os.Stderr, color, msg)
}
