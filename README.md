# Kuro

Kuro is a local-first meta-package manager for programming language environments, designed specifically for Arch Linux and CachyOS. 

Instead of acting as a traditional package manager that hosts and pulls binaries, Kuro acts as an intelligent wrapper and local registry over existing system tools like `pacman`, `paru`, and `yay`.

## Motivation

Managing programming languages on rolling-release distributions often leads to fragmented environments. A developer might install Go via pacman, a specific Java version via the AUR, and Node through a standalone script. 

Kuro exists to solve two core problems:
1. Centralized State: It maintains a localized, lightning-fast `registry.toml` of exactly what languages are installed, their specific versions, and how they were acquired. 
2. Standardized DX: It provides a uniform, visually clean CLI interface to install, remove, and update environments, abstracting away the underlying package managers.

## Installation

Kuro is written in Go. You can compile and install it globally using the native Go toolchain:

    git clone https://github.com/yourusername/kuro.git
    cd kuro
    go install

*(Ensure `~/go/bin` is in your `$PATH`)*

Alternatively, build the binary directly and move it to your system binaries:

    go build
    sudo mv kuro /usr/local/bin/

## Commands

Kuro operates identically to standard system managers, with strict constraints.

* `kuro install <language>`
  Resolves the target language, determines the correct upstream manager (pacman or AUR), safely executes the installation, and streams the output directly to the terminal UI.
  Example: `kuro install java8`

* `kuro remove <language>`
  Reverses the installation process, purging unused dependencies and dropping the record from the local state.
  Example: `kuro remove node`

* `kuro list`
  Instantly parses the local TOML registry and prints a structured table of all tracked environments without pinging the network.

* `kuro status`
  Queries the upstream repositories for all tracked languages, compares them against your local state, and provides an interactive prompt to selectively update outdated environments.

## Architecture

Kuro enforces strict separation of concerns to prevent command injection and state corruption:

* Router (`cmd/`): Built on Cobra. Strictly handles user intent and routing.
* State (`state/`): XDG-compliant, in-memory representations of the `~/.local/state/kuro/registry.toml` file.
* Resolver (`resolver/`): The logic layer. Translates user requests (e.g., `java23`) into real system packages and determines the safest tool (`pacman` vs `paru`) to acquire it.
* Runner (`runner/`): The execution layer. Spawns secure background processes without relying on shell expansion, piping standard output via Go's `io.Writer`.
* UI (`ui/`): A lightweight presentation layer ensuring raw terminal logs are formatted clearly without relying on heavy TUI libraries.

## License

MIT