# jot


## A Minimalist CLI Note-Taking Tool

jot is a plain text, CLI-first note-taking application designed for speed, focus, and UNIX-style composability.
It's inspired by tools like fzf, ripgrep, and navi, with a focus on workflow, not features.


### Install locally

```shell
git clone https://github.com/dalryan/jot
cd jot
go build
```

### Usage

```shell
dalryan@fedora ~/w/g/jot (main)> ./jot --help
Jot lets you quickly create, view, and organize plain-text notes from the command line.

Usage:
  jot [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  context     Manage the active context
  edit        Edit a note by ID
  help        Help about any command
  list        List existing notes
  quick       Capture a quick, timestamped note
  view        View a note by its ID

Flags:
  -h, --help   help for jot

Use "jot [command] --help" for more information about a command.

```

#### Example workflows

```shell
./jot quick "Idea for CLI tool" --tag jot,idea
./jot list --tag jot | fzf | xargs ./jot view
./jot edit 2a9f
./jot context set work
./jot quick "That prod issue again.." --tag k8s,bugs
./jot context clear
```


#### Use ripgrep to search

```shell
rg -l "work" $(./jot notes-path) | ./jot pipe
# or
alias jot-grep='rg -l "$1" $(./jot notes-path) | ./jot pipe'
jot-grep work
```


#### View timelines

```shell
./jot timeline --since 1h
./jot timeline --tag idea
./jot timeline --context work --tag infra
```

#### Export as json

```shell
./jot timeline --json | jq .
./jot list --json | jq '.[] | select(.Tags[] == "infra")'
rg -l "work" $(./jot notes-path) | ./jot pipe --json | jq .
```
