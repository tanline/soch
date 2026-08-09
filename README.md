# Soch

A lightweight CLI that scaffolds, manages, and centralizes notes across multiple repositories without duplicating files.

## Getting Started

### Prerequisites
- Go 1.26.5

### Setup
```shell
git clone https://github.com/tanline/soch.git
cd soch

make test
```

### Commands
```shell
make build  # Build the binary in bin/soch
make test   # Run tests
make clean  # Remove build artifacts
```

### Notes
```shell
soch init
soch new "Ideas"
soch add "A thought"
printf "A piped thought\n" | soch add
soch show
soch show 2026-08-08-ideas.md
```

`soch new` creates or selects a dated note. `soch add` appends a timestamped thought to the current note. `soch show` prints the current note or a named note.

### Manual testing

Build the binary, then use a disposable project, master-notes directory, and config file:

```yaml
# /tmp/soch-sandbox/config.yaml
master_notes_dir: /tmp/soch-sandbox/master-notes
```

```shell
make build

mkdir -p /tmp/soch-sandbox/project
mkdir -p /tmp/soch-sandbox/master-notes
cd /tmp/soch-sandbox/project

/path/to/soch/bin/soch --config /tmp/soch-sandbox/config.yaml init
/path/to/soch/bin/soch --config /tmp/soch-sandbox/config.yaml new "Ideas"
/path/to/soch/bin/soch --config /tmp/soch-sandbox/config.yaml add "A test thought"
/path/to/soch/bin/soch --config /tmp/soch-sandbox/config.yaml show
```

`--config` keeps this workflow separate from `~/.config/soch/config.yaml`.
