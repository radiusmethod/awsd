# awsd

AWS Profile Switcher in Go

Easily switch between AWS Profiles

<img src="assets/demo.gif" width="500">

## Requirements
min go 1.16

## Install

### Homebrew

```sh
brew tap pjauduiomv/awsd
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

Inspired by https://github.com/johnnyopao/awsp
