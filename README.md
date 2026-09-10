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

  gbrowse -commit [commit] opens the commit page.
  If commit is omitted, the current commit is opened.

  gbrowse -compare <COMPARE> [target] opens the comparison page between COMPARE and target.
  If target is omitted, the current commit is used.

Environment variables:
  GIT
    git command, default is git.

  DEBUG
    enable debug log if set.

Flags:
  -commit
    	open commit page
  -compare string
    	open compare page between the specified ref and target
  -print
    	only print generated url
```
