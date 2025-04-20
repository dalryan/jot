# jot


## A minimal and interoperable CLI Note-Taking tool

`jot` is a minimalist, Unix-style note-taking tool for the command line. It captures your thoughts fast, stores them in plain text, and plays well with everything else in your terminal toolbox.

> **Design Goals**
> - Work with stdin/stdout
> - Store notes as plaintext files
> - Be composable with `grep`, `awk`, `jq`, `fzf`, etc.
> - Require as little user configuration as possible
> - Clear organisation and context management

---

### Features

- **Capture instantly** from stdin or your editor
- **Composable** with Unix pipelines and scripting tools
- **Searchable** with `grep`, `jq`, `rg`, etc.
- **Plaintext** storage - portable and transparent
- **Minimal interface**, single config file, bring your own editor
- **Context-aware** note organization
- **Time-based** filtering and sorting
- **Tag-based** filtering and sorting
- **Template-driven** note creation
- **Markdown** note formatting

---

## Install

### Go install

```shell
go install github.com/dalryan/jot@v0.1.0
```

### Build from source

```shell
git clone https://github.com/dalryan/jot
cd jot
go build

# Optional: make `jot` available globally
alias jot=./jot
```

### Usage

```shell
dalryan@fedora ~/w/g/jot (main)> jot --help
Jot lets you quickly create, view, and organize plain-text notes from the command line.

Usage:
  jot [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  context     Manage the active context
  edit        Edit a note by ID
  help        Help about any command
  list        List existing notes
  new         Create a new note in your editor
  notes-path  Print the path to the notes directory
  pipe        Parse note file paths from stdin and display summaries
  quick       Capture a quick, timestamped note
  templates   Manage note templates
  timeline    Show notes in reverse chronological order
  today       Open or create today's daily note
  view        View a note by its ID

Flags:
  -h, --help   help for jot

Use "jot [command] --help" for more information about a command.
```

### Example workflows

The following examples show basic note-taking workflows using jot.

Example workflow:

```shell
# Quick note creation with tagging
jot quick "Idea for CLI tool" --tag jot,idea

# Context-based workflow management
jot context set work
jot quick "Production issue documentation" --tag k8s,bugs
jot context clear

# View all notes
jot list
# Edit a note by ID
jot edit <id>

# Time-based note filtering
jot timeline --since 1h
jot timeline --since 7d --tag idea
jot timeline --context work --tag k8s --since 1d
```

## Integration Capabilities

The strength of jot lies in its ability to compose — not replace — your existing toolset.

It operates effectively via standard streams, enabling seamless piping and filtering with other Unix utilities. 
Below are examples of real-world integration patterns.


### Integration with fzf (Command-Line Fuzzy Finder)
Select and view notes interactively:

```shell
# Tag-based note selection and viewing
jot list --tag idea | fzf | jot view

# Timeline-based note selection
jot timeline | fzf | jot view

# Interactive note selection for editing
jot list | fzf | jot edit
```

### Integration with ripgrep (rg)

Efficient content searching across notes:

```shell
# Content-based note filtering
rg -l "work" $(jot notes-path) | jot pipe

# Regular expression pattern matching
rg -l ".*bugs" $(jot notes-path) | jot pipe

# Combined content and tag filtering
rg -l "cli" $(jot notes-path) | jot pipe --tag idea
```

### Integration with jq (JSON Processor)

jot can output data in JSON format, enabling structured data processing with jq.

```shell
# JSON data structure visualization
jot timeline --json | jq .

# Tag-based JSON filtering
jot list --json | jq '.[] | select(.tags[] == "k8s")'

# Date-based filtering with field selection
jot timeline --json | jq '.[] | select(.created_at > "2023-01-01") | {id: .id, content: .content}'
```

### Integration with awk and sed (Text Processing)

Advanced text transformation and formatting:

```shell
# Text case transformation
jot list --json | jq -r '.[].Title' | sed 's/\(.*\)/\U\1/'

# Content volume analysis
find $(jot notes-path) -type f | xargs cat | wc -w
```


### Integration with find (File System Navigation)

Filter notes based on file attributes:

```shell
# Recent modification filtering
find $(jot notes-path) -type f -mtime -1 | jot pipe

# File size-based filtering
find $(jot notes-path) -type f -size +10k | jot pipe
```

### Integration with bat (Enhanced Output)

Syntax-highlighted and numbered output:

```shell
# Syntax highlighting for note content
jot view <id> | bat -l markdown

# Line number display for reference
jot view <id> | bat -l markdown --number
```

### Integration with pandoc (Document Converter)

Convert notes to other formats:

```shell
# HTML conversion for web publishing
jot view <id> | pandoc -f markdown -t html > note.html

# PDF document generation
jot view <id> | pandoc -f markdown -o note.pdf
```

---

jot helps you record your thoughts without breaking your flow — and without leaving your terminal.