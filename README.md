# frk

frk is summary generator for your GitHub activities.
(Golang port of [pepabo/furik](https://github.com/pepabo/furik))

## Installation

```
brew install winebarrel/frk/frk
```

## Usage

```
Usage: frk --token=STRING <command>

Flags:
  -h, --help            Show context-sensitive help.
      --version
      --token=STRING    GitHub token ($FRK_GITHUB_TOKEN)

Commands:
  activity --token=STRING
    show activity

  pulls --token=STRING
    show pull requests

Run "frk <command> --help" for more information on a command.
```
