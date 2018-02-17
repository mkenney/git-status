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

### New empty repository
```txt
$ git init && __git_status
local/master
```

### 1 untracked file, 1 total files
```txt
$ touch foo && __git_status
local/master ?1 #1
```

### 2nd untracked files, 2 total files
```txt
$ touch bar && __git_status
local/master ?2 #2
```

# 1 untracked file, 1 new file, 2 total files
```txt
$ git add foo && __git_status
local/master ?1 +1 #2
```

# 2 new files, 2 total files
```txt
$ git add bar && __git_status
local/master +2 #2
```

# 2 new files, 1 modified file, 1 modified file with unstaged changes, 2 total files
```txt
$ echo "baz" > foo && __git_status
local/master +2 𝚫1 ∴1 #2
```

# 2 new files, 2 modified file, 2 modified files with unstaged changes, 2 total files
```txt
$ echo "baz" > bar && __git_status
local/master +2 𝚫2 ∴2 #2
```

# 2 new files, 1 modified file, 1 modified file with unstaged changes, 2 total files

Because it's a newly tracked file, it sees it as a new file without changes once the changes are staged.

```txt
$ git add bar && __git_status
local/master +2 𝚫2 ∴2 #2
```

# clean working tree
```txt
$ git commit -am "commit" && __git_status
local/master
```

# 1 modified file, 1 modified file with unstaged changes, 1 total files
```txt
$ echo "baz2" >> bar && __git_status
local/master 𝚫1 ∴1 #1
```
