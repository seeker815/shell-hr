# SHR Command tool


## Introduction

Shell HR(shr) is a CLI tool build using Cobra(kubectl) to export a system's user information. The command will be able to export usernames, IDs, home directories, and shells as JSON or to stdout

 - By default, the command will display the information as JSON to stdout

 - Additionally, a file can be specified by using the -path flag

 - This command will not include information about system users (users with IDs under 1000). 

 - Here are the various ways the tools can be used:


## Examples

$ ./shr --help
```
shr command tool exports server's usernames, IDs, home directories as JSON

Usage:
  shr [flags]

Flags:
      
  -f, --format string   User data export format (default "json")
  -h, --help            help for shr
  -p, --path string     file path to export user data
  -t, --toggle          Help message for toggle
```

$ shr  -path=<path-file>.json
```
$ shr
[
  {
    "name": "sam",
    "id": 1002,
    "home": "/home/sam",
    "shell": "/bin/zsh"
  },
  {
    "name": "ubuntu",
    "id": 1003,
    "home": "/home/ubuntu",
    "shell": "/bin/bash"
  },
]
```
