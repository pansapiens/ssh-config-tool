#!/usr/bin/env python
import os
import re
import argparse
import sys

from ssh_config_tool import (
    ssh_command_to_config_entry,
    split_ssh_config,
    backup_ssh_config,
    parse_ssh_config,
)


def translate(args):
    ssh_command = " ".join(args.ssh_args)
    config_entry = ssh_command_to_config_entry(ssh_command)

    ssh_config_d = os.path.expanduser("~/.ssh/config.d")
    os.makedirs(ssh_config_d, exist_ok=True)

    host = re.search(r"^Host\s+(\S+)", config_entry, re.MULTILINE).group(1)
    host_file = os.path.join(ssh_config_d, f"{host}.conf")

    with open(host_file, "w") as f:
        f.write(config_entry)

    print(f"Generated {host_file}:")

    with open(host_file, "r") as f:
        print(f.read())


def main():
    parser = argparse.ArgumentParser(description="SSH config management tool")
    subparsers = parser.add_subparsers()

    split_parser = subparsers.add_parser(
        "split", help="Split SSH config into separate files"
    )
    split_parser.set_defaults(func=split_ssh_config)

    translate_parser = subparsers.add_parser(
        "translate", help="Translate SSH command to SSH config entry"
    )
    translate_parser.add_argument(
        "ssh_args", nargs=argparse.REMAINDER, help="SSH command arguments"
    )
    translate_parser.set_defaults(func=translate)

    args = parser.parse_args()

    if "func" in args:
        args.func(args)
    else:
        parser.print_help()
        sys.exit(1)


if __name__ == "__main__":
    main()
