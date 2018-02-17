# git-status

This is a simple git status line for your shell prompt. Defines a function called `__git_status()` that returns a string indicating the current state of your local repository. Displays:

* origin/local/detached/tag indicator
* commit name
* number of commits behind origin
* number of commits ahead of origin
* number of untracked files
* number of deleted files
* number of added files
* number of renamed files
* number of modified files
* number of files with unstaged changes
* total number of files

## examples

### New repository
```txt
$ git init && __git_status
local/master
```

### Untracked files
```txt
$ touch foo && __git_status
local/master ?1 #1
```
```txt
$ touch bar && __git_status
local/master ?2 #2
```

### Added files
```txt
$ git add foo && __git_status
local/master ?1 +1 #2
```
```txt
$ git add bar && __git_status
local/master +2 #2
```

### Modified files
```txt
$ echo "baz" > foo && __git_status
local/master +2 𝚫1 ⸮1 #3
```
```txt
$ echo "baz" > bar && __git_status
local/master +2 #2
```
