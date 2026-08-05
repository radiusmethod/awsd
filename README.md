# awsd - AWS Profile Switcher in Go

---

<img src="assets/awsd.png" width="200">

awsd is a command-line utility that allows you to easily switch between AWS Profiles.

<img src="assets/demo.gif" width="500">

## Table of Contents

- [Installation](#installation)
    - [Homebrew](#homebrew)
    - [Prebuilt binary](#prebuilt-binary)
    - [Makefile](#makefile)
    - [To Finish Installation](#to-finish-installation)
    - [Upgrading](#upgrading)
    - [Upgrading from pre-v0.4.0](#upgrading-from-pre-v040)
- [Usage](#usage)
    - [Switching AWS Profiles](#switching-aws-profiles)
    - [Switching AWS Regions](#switching-aws-regions)
    - [Show your AWS Profile in your shell prompt](#show-your-aws-profile-in-your-shell-prompt)
    - [Add autocompletion](#add-autocompletion)
- [Why a shell function?](#why-a-shell-function)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Homebrew

```sh
brew tap radiusmethod/awsd
brew install awsd
```

### Prebuilt binary

Grab the archive for your platform from the
[latest release](https://github.com/radiusmethod/awsd/releases/latest), then put `awsd` somewhere
on your `PATH`. macOS, Linux, and Windows on amd64 and arm64.

### Makefile

Builds from source, so this one needs [Go](https://golang.org/dl/) installed.

```sh
make install
```

### To Finish Installation
Add one line to your shell's startup file, then open a new terminal or source that file.

**zsh** (`~/.zshrc`):
```sh
eval "$(awsd init zsh)"
```

**bash** (`~/.bashrc` or `~/.bash_profile`):
```sh
eval "$(awsd init bash)"
```

**fish** (`~/.config/fish/config.fish`):
```fish
awsd init fish | source
```

**PowerShell** (`$PROFILE`):
```powershell
awsd init powershell | Out-String | Invoke-Expression
```

Ex. `echo 'eval "$(awsd init zsh)"' >> ~/.zshrc`

That one line defines the `awsd` shell function, sets up tab completion, and applies the profile
and region you last selected to every new shell. Nothing else to configure.

### Upgrading
Upgrading consists of just doing a brew update and brew upgrade.

```sh
brew update && brew upgrade radiusmethod/awsd/awsd
```

### Upgrading from pre-v0.4.0
**v0.4.0 is a breaking change.** awsd now installs one binary named `awsd`. The old
`_awsd_prompt` binary, the `_awsd` wrapper script, and the `_awsd_autocomplete` completion script
are all gone, and so is the alias-based setup.

Delete whatever of this you have in your shell config:

```sh
alias awsd="source _awsd"        # removed in v0.4.0
source _awsd_autocomplete        # removed in v0.4.0
if [ -f ~/.awsd ]; then ...      # removed in v0.3.0
```

and replace it with the single line for your shell from
[To Finish Installation](#to-finish-installation) above.

Then clear out the old files, which a package manager will not remove for you if you ever ran
`make install` by hand:

```sh
rm -f /usr/local/bin/_awsd_prompt /usr/local/bin/_awsd /usr/local/bin/_awsd_autocomplete
type -a _awsd_prompt   # should print nothing
```

Two things that bite during the upgrade:

- **Remove the old alias.** In zsh an alias shadows a function of the same name, so leaving
  `alias awsd="source _awsd"` in place means the new `awsd` function never gets used. If the alias
  is defined *before* the `eval` line, the eval fails outright with
  `defining function based on alias 'awsd'`.
- **Put the `eval` line after any `PATH` changes** that point at your awsd install. It runs `awsd`
  at startup, so if an older copy is earlier in `PATH` at that moment you can get
  `(eval):1: bad pattern: ^[[0`, which is an old binary printing a colored warning that your shell
  then tries to eval. `type -a awsd` shows you every copy.

Your `~/.awsd` file carries over untouched. Profile and region selections survive the upgrade.

## Usage

### Switching AWS Profiles

It is possible to shortcut the menu selection by passing the profile name you want to switch to as an argument.

```bash
> awsd work
Profile work set.
```

To switch between different profiles files using the menu, use the following command:

```bash
awsd
```

This command will display a list of available profiles files in your `~/.aws/config` file or from `AWS_CONFIG_FILE`
if you have that set. It expects for you to have named profiles in your AWS config file. Select the one you want to use.

### Switching AWS Regions

You can also switch your active AWS region. The interactive picker fuzzy-matches the same way the profile picker does.

```bash
> awsd set region us-east-1
Region us-east-1 set.

> awsd set region          # interactive picker
> awsd list regions        # list known regions
> awsd unset region        # clear the active region
```

Setting a region exports `AWS_REGION` and `AWS_DEFAULT_REGION` in the calling shell. Profile and region are independent — `awsd set profile` does not change your region, and vice versa.

Your selection persists across new terminal windows automatically, since `awsd init` applies
whatever is in `~/.awsd` when each shell starts.

### Show your AWS Profile in your shell prompt
For better visibility into what your shell is set to it can be helpful to configure your prompt to show the value of the env variable `AWS_PROFILE`.

<img src="assets/screenshot.png" width="700">

Here's a sample of my zsh prompt config using oh-my-zsh themes

```sh
# AWS info
local aws_info='$(aws_prof)'
function aws_prof {
  local profile="${AWS_PROFILE:=}"
  echo -n "%{$fg_bold[blue]%}aws:(%{$fg[cyan]%}${profile}%{$fg_bold[blue]%})%{$reset_color%} "
}
```

```sh
PROMPT='OTHER_PROMPT_STUFF $(aws_info)'
```

### Add autocompletion
Tab completion comes with `awsd init`. It completes profile names on `awsd <TAB>`, the
`set`/`unset`/`list` subcommands, and their arguments, so `awsd set region <TAB>` lists regions
and `awsd set profile <TAB>` lists profiles.

## Why a shell function?

`awsd init` generates a shell function rather than shipping a plain binary, because a child
process cannot change its parent shell's environment. Anything that sets `AWS_PROFILE` for your
current shell has to run *in* that shell.

So the binary does the picking and writes your choice to `~/.awsd`, and the generated function
asks it for the matching shell code and evals that:

```sh
awsd() {
  command awsd "$@" || return
  eval "$(command awsd shellenv bash)"
}
```

The function and the binary share the name `awsd`. That works because `command` skips functions
and aliases and runs the executable from `PATH`. PowerShell's `&` operator does not do this, so
the generated PowerShell integration resolves the binary path up front with `Get-Command` instead.

You can see exactly what gets eval'd at any time:

```sh
awsd init zsh      # the whole integration
awsd shellenv zsh  # just the exports for the current selection
```

## Contributing

If you encounter any issues or have suggestions for improvements, please open an issue or create a pull request on [GitHub](https://github.com/radiusmethod/awsd).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


Inspired by https://github.com/johnnyopao/awsp
