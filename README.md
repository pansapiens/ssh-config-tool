# sshcfgtool

This tool is designed to tame the handling of `~/.ssh/config` files. 

It can be used to split the config file into multiple files in `~/.ssh/config.d`, and generate new config files from a provided `ssh` command line example.

```
Usage: sshcfgtool [subcommand] [options]

Subcommands:
  split         Split SSH config into separate files
  translate     Translate SSH command to SSH config entry

Global flags:
  -c, -config string
        Path to the SSH config file
  -n, -dry-run
        Print intentions but don't actually change any files
  -force
        Force overwriting existing files in split subcommand
```

## Splitting

To split your existing ~/.ssh/config into separate files named {host}.conf located in ~/.ssh/config.d:

```bash
sshcfgtool split
```

The existing `~/.ssh/config` file will be renamed to `~/.ssh/config.backup` and a new config file generated containing `Include /home/perry/.ssh/config.d/*`.

## Translating an SSH command line to a config file entry

To convert an `ssh` command line to a config file entry, write your usual command after the `translate` subcommand:

```bash
sshcfgtool translate ssh -i ~/.ssh/somekey -A -X -p 2222 user@hostname | tee ~/.ssh/config.d/hostname.conf
```

_For better or worse, the code for this project was primarily written with the guidance of OpenAI's ChatGPT (GPT-4). It will probably be a very normal thing in a few years. Both Go and Python implementations are available. The Go version is currently the most up-to-date and includes unit tests. However, the final decision on which version will be the primary implementation has not been made._

## Development

### Test
```bash
go test -v
```

### Run
```bash
go run ssh_config_tool.go
```

### Build

```bash
go build -o sshcfgtool ssh_config_tool.go
```

## Dependencies

Dependencies are managed using Go modules, with dependencies already defined in `go.mod`.

You don't need to run the commands below, since `go run` and `go build` will fetch dependencies as required.
For reference this is how the module was initialized and dependencies were installed:

```bash
go mod init ssh-config-tool

go get github.com/pmezard/go-difflib/difflib
go get github.com/google/shlex
```