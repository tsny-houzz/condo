# Condo

`condo` is a command line utility that allows you to create `codespaces`. 

`codespaces` are isolated kubernetes testing environments.

## Prerequisites

- Have your email in your `~/.gitconfig`
- Your k8s config must be in the `stg-main-eks` cluster

## Usage

```
NAME:
   condo - CLI tool to manage namespaces and applications in codespaces

USAGE:
   condo [global options] command [command options]

COMMANDS:
   whoami   Show the user's email
   ns       Namespace-related commands
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

------------------------------------------------------------------------
❯ condo ns     

NAME:
   condo ns - Codespace-related commands

USAGE:
   condo ns command [command options]

COMMANDS:
   list     List all codespaces
   create   Create a codespace
   desc     Describe a codespace
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

## Example Commands

```bash
condo ns list --owner yj@houzz.com
```

## Building

All deps are from public go repos.

`go install` to create a binary to your `GOBIN` or `go build` and move the binary to somewhere in your `PATH` to use from anywhere.

