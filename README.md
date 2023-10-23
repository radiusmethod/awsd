# awsd

AWS Profile Switcher in Go

Easily switch between AWS Profiles

<img src="assets/demo.gif" width="500">

It is possible to short cut the menu selection by passing
the profile name you want to switch to as an argument.

```sh
> awsd work
Profile work set.
```


## Requirements
min go 1.16

## Install

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

## Usage
```sh
awsd
```

## Persist Profile across new shells
To persist the set profile when you open new terminal windows, you can add the following to your bash profile or zshrc.

```bash
export AWS_PROFILE=$(cat ~/.awsd)
```

## Show your AWS Profile in your shell prompt
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

## Add autocompletion
You can add autocompletion when passing profile as argument by creating a script with the following. I put it in 
`~/bin/awsd_autocomplete.sh`, then source that script and add to your bash profile or zshrc file.
`source ~/bin/awsd_autocomplete.sh`

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

Inspired by https://github.com/johnnyopao/awsp
