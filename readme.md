# Condo

**NOTE**: THIS IS A WIP EXPERIMENT FOR THE FUTURE OF CODEPATH/JUKWAA DEBUG (jkdebug cookie)

`condo` is a command line utility that allows you to create `codespaces`. 

`codespaces` are isolated kubernetes testing environments.


## Prerequisites

- Setup [Teleport](https://engwiki.houzz.tools/doc/k8s-access-using-teleport-r4JAMoor01)
- Have your email in your `~/.gitconfig`
- k8s
   - Your k8s config must be in the `stg-main-eks` cluster
   - You must have permissions to create namespaces and operate in the `stg-main-eks` cluster

## Usage

```
NAME:
   condo - CLI tool to manage codespaces

USAGE:
   condo [global options] command [command options]

COMMANDS:
   whoami   Show the user's email
   cs       Codespace-related commands
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

------------------------------------------------------------------------
❯ condo cs     

NAME:
   condo cs - Codespace-related commands

USAGE:
   condo cs command [command options]

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
# Lists all codespaces owned by yj
condo cs list --owner yj@houzz.com

# Creates a codespace called `tsny-testing-cs`
condo cs create tsny-testing-cs

# Installs an applicationo to the `tsny-testing-cs` codespace
condo cs install tsny-testing-cs prismic-cms-code:feat-seo2-2541_2024_10_31__18_01_13_3b3edf0f44 
```

## Building

All deps are from public go repos.

`go install` to create a binary to your `GOBIN` or `go build` and move the binary to somewhere in your `PATH` to use from anywhere.

