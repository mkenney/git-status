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
```sh
origin/master
```



`[origin]/[commit] [commits behind origin] [commits ahead of origin] [untracked files]
