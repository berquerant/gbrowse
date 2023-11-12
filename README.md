# gbrowse

```
❯ gbrowse -h
gbrowse - Open the repo in the browser

Usage:
  gbrowse [flags] [target]

  The target is PATH or FILE:LINUM.
  gbrowse PATH opens the PATH of the repo.
  gbrowse FILE:LINUM opens the line LINUM of the FILE of the repo.
  gbrowse opens the directory of the repo.

Config:

  {
    "phases": [
      PHASE, ...
    ],
    "defs": [
      {
        "id": ID,
        "cmd": ["command", ...]
      }, ...
    ]
  }

phases determines the search order for ref (commit, branch, tag).
PHASE is branch, default_branch, tag, commit or id in def.
defs is custom phases, cmd should return a string like commit hash, for example,

  {
    "phases": ["echo-master"],
    "defs": [{"id": "echo-master", "cmd": ["echo", "master"]}]
  }

sets ref to "master".

If all searches fail, search commit.

Environment variables:
  GBROWSE_GIT
    git command, default is git.

  GBROWSE_DEBUG
    enable debug log if set.

  GBROWSE_CONFIG
    config file or string.
    -config overwrites this.

Flags:
  -config string
        config or file
  -print
        only print generated url
```
