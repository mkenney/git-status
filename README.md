# git-status

This is a simple, clean, informative `git` status line for your `bash` shell prompt. The `git-status` script defines a function called, creatively, `__git_status()` that returns a string indicating the current state of your local repository. The function returns a string describing:

* `origin`/`local`/`tag`/`detached` origin indicator
* branch name/tag name/commit hash position indicator
* total commits behind origin: `↓n`
* total commits ahead of origin: `↑n`
* total untracked files: `…n`
* total deleted files: `✖n`
* total added files: `✚n`
* total renamed files: `↪n`
* total staged files: `✔n`
* total unstaged files: `✎n`
* total number of files: `#n`

A complex set of changes containing all these elements might produce a status line that looks something like:

```txt
origin/some-feature/mybranch ↓2 ↑2 ✖1 ✚1 ↪1 ✔2 ✎1 …1 #5
```

though that doesn't really happen much. I rarely have more than 1 - 3 status indicators showing at any given time. ymmv.

## usage

Running all the `git` commands sequentially in `bash` is a bit slow, sometimes pushing a full second for complex changes or a large number of files, so I reimplemented it in `go`. The shell script will still fallback to the `bash` version if the binaries aren't found.

The performance of the `go` version is limited to the speed of the slowest `git` command, generally `git diff --name-only` to list the changed files. Usually 100 - 200 milliseconds.

To enable the `go` version, add [`git-status-darwin-amd64`](https://github.com/mkenney/git-status/blob/go/bin/git-status-darwin-amd64) and/or [`git-status-linux-amd64`](https://github.com/mkenney/git-status/blob/go/bin/git-status-linux-amd64) to your path. `git-status` will detect the correct platform automatically.

There's a [`Makefile`](https://github.com/mkenney/git-status/blob/go/Makefile) now... so adding more architectures is easy enough.

## examples

#### detached head
```txt
$ git checkout <some hash> && __git_status
detached/20528e7ad4
```

#### tagged commit
```txt
$ git tag v0.0.1 && __git_status
tag/v0.0.1
```

#### New empty repository
```txt
$ git init && __git_status
local/master
```

#### 1 untracked file, 1 total files
```txt
$ touch foo && __git_status
local/master …1 #1
```

#### 2nd untracked files, 2 total files
```txt
$ touch bar && __git_status
local/master …2 #2
```

#### 1 untracked file, 1 new file, 2 total files
```txt
$ git add foo && __git_status
local/master …1 +1 #2
```

#### 2 new files, 2 total files
```txt
$ git add bar && __git_status
local/master +2 #2
```

#### 2 new files, 1 modified file, 1 file with unstaged changes, 2 total files
```txt
$ echo "baz" > foo && __git_status
local/master +2 ✔1 ✎1 #2
```

#### 2 new files, 2 modified file, 2 files with unstaged changes, 2 total files
```txt
$ echo "baz" > bar && __git_status
local/master +2 ✔2 ✎2 #2
```

#### 2 new files, 1 modified file, 1 file with unstaged changes, 2 total files

Because it's a newly tracked file, it sees it as a new file without changes once the changes are staged.

```txt
$ git add bar && __git_status
local/master +2 ✔2 ✎2 #2
```

#### clean working tree
```txt
$ git commit -am "commit" && __git_status
local/master
```

#### 1 renamed file, 1 total files
```txt
$ git mv bar baz && __git_status && git reset --hard
local/master ↪1 #1
```

#### 1 modified file, 1 file with unstaged changes, 1 total files
```txt
$ echo "baz2" >> bar && __git_status
local/master ✔1 ✎1 #1
```

#### 2 modified files, 2 files with unstaged changes, 2 total files
```txt
$ echo "baz" >> foo && __git_status
local/master ✔2 ✎2 #2
```

#### 2 modified files, 1 file with unstaged changes, 2 total files
```txt
$ git add foo && __git_status
local/master ✔2 ✎1 #2
```

#### 1 untracked file, 2 modified files, 1 file with unstaged changes, 3 total files
```txt
$ touch baz && __git_status
local/master …1 ✔2 ✎1 #3
```

#### 2 untracked files, 2 modified files, 1 file with unstaged changes, 4 total files
```txt
$ touch 00ntz && __git_status
local/master …2 ✔2 ✎1 #4
```

#### 1 untracked file, 1 new file, 2 modified files, 1 file with unstaged changes, 4 total files
```txt
$ git add baz && __git_status
local/master …1 +1 ✔2 ✎1 #4
```

#### 1 untracked file, 1 deleted file, 1 new file, 1 modified file, 1 file with unstaged changes, 4 total files
```txt
$ git rm -f foo && __git_status
local/master …1 ×1 +1 ✔1 ✎1 #4
```

#### clean working tree
```txt
$ git commit -am "commit" && __git_status
local/master
```

#### branch origin set
```txt
$ git remote add origin https://github.com/user/repo.git && git push -u origin master && __git_status
origin/master
```

#### origin 2 commits ahead of local branch
```txt
$ ... && __git_status
origin/master ↓2
```

#### local branch 2 commits ahead of origin
```txt
$ ... && __git_status
origin/master ↑2
```

#### local branch contains 2 commits origin doens't have and origin contains 2 commits local branch doesn't have...
```txt
$ ... && __git_status
origin/master ↓2 ↑2
```

#### All together
##### local ahead 2 commits, origin ahead 2 commits, 1 untracked file, 1 deleted file, 1 new file, 2 modified files, 1 renamed file, 1 file with unstaged changes, 5 total files
```txt
$ ... && __git_status
origin/master ↓2 ↑2 …1 ×1 +1 ↪1 ✔2 ✎1 #5
```
