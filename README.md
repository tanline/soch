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
```

`soch new` creates or selects a dated collection. `soch add` appends a timestamped thought to the current collection.
