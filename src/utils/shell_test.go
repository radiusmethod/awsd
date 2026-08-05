package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseShell(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Shell
		expectError bool
	}{
		{name: "bash", input: "bash", expected: Bash},
		{name: "zsh", input: "zsh", expected: Zsh},
		{name: "fish", input: "fish", expected: Fish},
		{name: "powershell", input: "powershell", expected: PowerShell},
		{name: "pwsh alias", input: "pwsh", expected: PowerShell},
		{name: "case insensitive", input: "PowerShell", expected: PowerShell},
		{name: "surrounding space", input: "  zsh ", expected: Zsh},
		{name: "unknown shell", input: "csh", expectError: true},
		{name: "empty", input: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := ParseShell(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, shell)
		})
	}
}

func TestShellEnv(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		shell    Shell
		expected string
	}{
		{
			name:     "Profile only, posix",
			state:    State{Profile: "dev"},
			shell:    Bash,
			expected: "export AWS_PROFILE='dev'\n",
		},
		{
			name:     "Profile and region, posix",
			state:    State{Profile: "dev", Region: "us-east-1", RegionSet: true},
			shell:    Zsh,
			expected: "export AWS_PROFILE='dev'\nexport AWS_REGION='us-east-1'\nexport AWS_DEFAULT_REGION='us-east-1'\n",
		},
		{
			name:     "Empty profile unsets, posix",
			state:    State{},
			shell:    Bash,
			expected: "unset AWS_PROFILE\n",
		},
		{
			name:     "Explicit empty region unsets both, posix",
			state:    State{Profile: "dev", RegionSet: true},
			shell:    Bash,
			expected: "export AWS_PROFILE='dev'\nunset AWS_REGION AWS_DEFAULT_REGION\n",
		},
		{
			name:     "Profile and region, fish",
			state:    State{Profile: "dev", Region: "us-east-1", RegionSet: true},
			shell:    Fish,
			expected: "set -gx AWS_PROFILE 'dev'\nset -gx AWS_REGION 'us-east-1'\nset -gx AWS_DEFAULT_REGION 'us-east-1'\n",
		},
		{
			name:     "Empty profile and empty region, fish",
			state:    State{RegionSet: true},
			shell:    Fish,
			expected: "set -e AWS_PROFILE\nset -e AWS_REGION AWS_DEFAULT_REGION\n",
		},
		{
			name:     "Profile and region, powershell",
			state:    State{Profile: "dev", Region: "us-east-1", RegionSet: true},
			shell:    PowerShell,
			expected: "$env:AWS_PROFILE = 'dev'\n$env:AWS_REGION = 'us-east-1'\n$env:AWS_DEFAULT_REGION = 'us-east-1'\n",
		},
		{
			name:     "Empty profile and empty region, powershell",
			state:    State{RegionSet: true},
			shell:    PowerShell,
			expected: "$env:AWS_PROFILE = $null\n$env:AWS_REGION = $null\n$env:AWS_DEFAULT_REGION = $null\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ShellEnv(tt.state, tt.shell))
		})
	}
}

// Profile names come from ~/.aws/config, so they can contain anything a section
// header allows. Unquoted output would let that run as shell code.
func TestShellEnvQuoting(t *testing.T) {
	tests := []struct {
		name     string
		shell    Shell
		profile  string
		expected string
	}{
		{
			name:     "Space, posix",
			shell:    Bash,
			profile:  "my profile",
			expected: "export AWS_PROFILE='my profile'\n",
		},
		{
			name:     "Single quote, posix",
			shell:    Bash,
			profile:  "we'ird",
			expected: `export AWS_PROFILE='we'\''ird'` + "\n",
		},
		{
			name:     "Command substitution stays literal, posix",
			shell:    Bash,
			profile:  "$(whoami)`id`",
			expected: "export AWS_PROFILE='$(whoami)`id`'\n",
		},
		{
			name:     "Single quote, fish",
			shell:    Fish,
			profile:  "we'ird",
			expected: `set -gx AWS_PROFILE 'we\'ird'` + "\n",
		},
		{
			name:     "Backslash, fish",
			shell:    Fish,
			profile:  `back\slash`,
			expected: `set -gx AWS_PROFILE 'back\\slash'` + "\n",
		},
		{
			name:     "Single quote, powershell",
			shell:    PowerShell,
			profile:  "we'ird",
			expected: "$env:AWS_PROFILE = 'we''ird'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShellEnv(State{Profile: tt.profile}, tt.shell)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestInitScript(t *testing.T) {
	tests := []struct {
		name     string
		shell    Shell
		contains []string
	}{
		{
			name:     "bash",
			shell:    Bash,
			contains: []string{"awsd() {", `eval "$(command awsd shellenv bash)"`, "complete -o nospace -F _awsd_completion awsd"},
		},
		{
			name:     "zsh",
			shell:    Zsh,
			contains: []string{"awsd() {", `eval "$(command awsd shellenv zsh)"`, "bashcompinit"},
		},
		{
			name:     "fish",
			shell:    Fish,
			contains: []string{"function awsd", "command awsd shellenv fish | source", "complete -c awsd"},
		},
		{
			name:     "powershell",
			shell:    PowerShell,
			contains: []string{"function awsd", "& $global:AwsdBin shellenv powershell | Out-String | Invoke-Expression", "Register-ArgumentCompleter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := InitScript(tt.shell)
			assert.NotEmpty(t, script)
			for _, want := range tt.contains {
				assert.Contains(t, script, want)
			}
		})
	}
}

// bashcompinit only provides `complete` for zsh; bash must not carry it.
func TestInitScriptBashHasNoCompinit(t *testing.T) {
	assert.NotContains(t, InitScript(Bash), "bashcompinit")
}

func TestInitScriptUnknownShell(t *testing.T) {
	assert.Empty(t, InitScript(Shell("csh")))
}

// Every accepted shell has to produce a script, or `awsd init <shell>` would
// validate the argument and then print nothing.
func TestInitScriptCoversAcceptedShells(t *testing.T) {
	for _, name := range AcceptedShells {
		shell, err := ParseShell(name)
		assert.NoError(t, err)
		assert.NotEmpty(t, InitScript(shell), "no init script for %s", name)
		assert.True(t, strings.HasSuffix(InitScript(shell), "\n"), "%s script must end in a newline", name)
	}
}
