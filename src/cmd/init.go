package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/radiusmethod/awsd/src/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <shell>",
	Short: "Print the shell integration for the given shell.",
	Long: "Print the shell integration for the given shell: the awsd function, completion, " +
		"and the hook that applies the profile/region in ~/.awsd to new shells.\n\n" +
		"awsd has to run inside your shell to export AWS_PROFILE, because a child process " +
		"cannot change its parent's environment. Eval this from your rc file:\n\n" +
		"  bash/zsh:   eval \"$(awsd init zsh)\"\n" +
		"  fish:       awsd init fish | source\n" +
		"  PowerShell: awsd init powershell | Out-String | Invoke-Expression",
	ValidArgs: utils.AcceptedShells,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		shell, err := utils.ParseShell(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(utils.InitScript(shell))
	},
}

var shellenvCmd = &cobra.Command{
	Use:   "shellenv [shell]",
	Short: "Print the shell code that applies ~/.awsd to the current shell.",
	Long: "Print the export/unset statements for the active profile and region. " +
		"Used by the function that `awsd init` generates; you should not need to call it directly. " +
		"Defaults to POSIX (bash/zsh) syntax.",
	ValidArgs: utils.AcceptedShells,
	Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		shell := utils.Bash
		if len(args) == 1 {
			s, err := utils.ParseShell(args[0])
			if err != nil {
				log.Fatal(err)
			}
			shell = s
		}
		homeDir, err := utils.GetHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		state, err := utils.ReadState(homeDir)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(utils.ShellEnv(state, shell))
	},
}

func init() {
	initCmd.Example = fmt.Sprintf("  eval \"$(awsd init zsh)\"\n\n  # supported: %s", strings.Join(utils.SupportedShells, ", "))
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(shellenvCmd)
}
