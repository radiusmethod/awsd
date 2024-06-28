# awsd - AWS Profile Switcher in Go

---

<img src="assets/awsd.png" width="200">

awsd is a command-line utility that allows you to easily switch between AWS Profiles.

<img src="assets/demo.gif" width="500">

## Table of Contents

- [Installation](#installation)
    - [Homebrew](#homebrew)
    - [Makefile](#makefile)
    - [To Finish Installation](#to-finish-installation)
    - [Upgrading](#upgrading)
- [Usage](#usage)
    - [Switching AWS Profiles](#switching-aws-profiles)
    - [Persist Profile across new shells](#persist-profile-across-new-shells)
    - [Show your AWS Profile in your shell prompt](#show-your-aws-profile-in-your-shell-prompt)
    - [Add autocompletion](#add-autocompletion)
    - [TL;DR (full config example)](#tldr-full-config-example)
- [Contributing](#contributing)
- [License](#license)

## Installation

Make sure you have Go installed. You can download it from [here](https://golang.org/dl/).

### Homebrew

```sh
brew tap radiusmethod/awsd
brew install awsd
```

### Makefile

```sh
make install
```

### To Finish Installation
Add the following to your bash profile or zshrc then open new terminal or source that file

```sh
alias awsd="source _awsd"
```

Ex. `echo 'alias awsd="source _awsd"' >> ~/.zshrc`

### Upgrading
Upgrading consists of just doing a brew update and brew upgrade.

```sh
brew update && brew upgrade radiusmethod/awsd/awsd
```

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

### Persist Profile across new shells
To persist the set profile when you open new terminal windows, you can add the following to your bash profile or zshrc.

```bash
export AWS_PROFILE=$(cat ~/.awsd)
```

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
You can add autocompletion when passing config as argument by adding the following to your bash profile or zshrc file.
`source _awsd_autocomplete`

```bash
[ "$BASH_VERSION" ] && AWSD_CMD="awsd" || AWSD_CMD="_awsd"
_awsd_completion() {
    local cur=${COMP_WORDS[COMP_CWORD]}
    local suggestions=$(awsd list)
    COMPREPLY=($(compgen -W "$suggestions" -- $cur))
    return 0
}
complete -o nospace -F _awsd_completion "${AWSD_CMD}"
```

Now you can do `awsd my-p` and hit tab and if you had a profile `my-profile` it would autocomplete and find it.

### TL;DR (full config example)
```bash
alias awsd="source _awsd"
source _awsd_autocomplete
export AWS_PROFILE=$(cat ~/.awsd)
```

## Contributing

If you encounter any issues or have suggestions for improvements, please open an issue or create a pull request on [GitHub](https://github.com/radiusmethod/awsd).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


Inspired by https://github.com/johnnyopao/awsp
