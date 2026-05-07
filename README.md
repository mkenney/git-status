# git-status

A simple, clean, informative `git` status line for your shell prompt. The `git-status` script defines a function called `__git_status()` that returns a string describing the current state of your local repository.

Output format:

```txt
⎇ <origin>: <position> […N] [＊N] [↓N] [↑N] [✖ N] [✚ N] [↪ N] [✔ N] [✎ N]
```

Indicators, in output order:

| Indicator | Meaning |
|-----------|---------|
| `origin:` | `origin`, `local`, `tag`, or `detached` |
| `position` | branch name, tag name, or short commit hash |
| `…N` | untracked files |
| `＊N` | stashed changes |
| `↓N` | commits behind origin |
| `↑N` | commits ahead of origin |
| `✖ N` | deleted files |
| `✚ N` | added (newly staged) files |
| `↪ N` | renamed files |
| `✔ N` | staged modifications (clean working tree) |
| `✎ N` | files with unstaged changes |

A complex status might look like:

```txt
⎇ origin: main …1 ＊1 ↓2 ↑2 ✖ 1 ✚ 1 ↪ 1 ✔ 2 ✎ 1
```

though in practice you'll rarely see more than 1–3 indicators at a time.

## Usage

Source the [`git-status`](https://github.com/mkenney/git-status/blob/master/git-status) file in your shell to define the `__git_status` function, then call it from your prompt:

```bash
source /path/to/git-status
# In your PS1 or precmd hook:
PS1='$([ -d .git ] && __git_status) $ '
```

## Performance

The shell script fallback uses a single `git status --porcelain=v2 --branch` call (plus `git stash list`) to collect all state, replacing the ~10 sequential git subprocess calls of the naive approach.

The Go binary is faster still — it fans out 6 git commands in parallel goroutines and launches the `rev-list` ahead/behind query as soon as the hash and upstream data are available, without waiting for the other commands to complete. Typical wall time is 100–200 ms.

### Binary installation

Pre-built binaries are available in the `bin/` directory. Add the appropriate binary for your platform to your `PATH`:

| File | Platform |
|------|----------|
| `git-status-darwin-arm64` | macOS Apple Silicon |
| `git-status-darwin-amd64` | macOS Intel (also works on Apple Silicon via Rosetta) |
| `git-status-linux-arm64` | Linux aarch64 (AWS Graviton, RPi 4+) |
| `git-status-linux-amd64` | Linux x86-64 |
| `git-status-linux-armv7` | Linux ARMv7 (RPi 2/3) |

The shell script detects the correct binary automatically based on `uname` and `uname -m`, preferring the native-arch binary when available.

### Building from source

```bash
# All architectures
make build

# Single architecture
make build-darwin-arm64
make build-darwin-amd64
make build-linux-arm64
make build-linux-amd64
make build-linux-armv7
```

## Examples

#### Clean working tree, local branch

```txt
⎇ local: master
```

#### Tracking origin, clean

```txt
⎇ origin: main
```

#### Detached HEAD

```txt
⎇ detached: 20528e7ad4
```

#### Tagged commit (detached HEAD at a tag)

```txt
⎇ tag: v1.2.3 (20528e7ad4)
```

#### 1 untracked file

```txt
⎇ local: master …1
```

#### 2 staged new files

```txt
⎇ local: master ✚ 2
```

#### 2 staged new files, 1 with unstaged changes

```txt
⎇ local: master ✚ 2 ✔ 1 ✎ 1
```

#### 1 renamed file

```txt
⎇ local: master ↪ 1
```

#### Origin 2 commits ahead of local

```txt
⎇ origin: main ↓2
```

#### Local 2 commits ahead of origin

```txt
⎇ origin: main ↑2
```

#### Diverged: 2 commits in each direction

```txt
⎇ origin: main ↓2 ↑2
```

#### 1 stashed change

```txt
⎇ origin: main ＊1
```

#### All indicators at once

```txt
⎇ origin: main …1 ＊1 ↓2 ↑2 ✖ 1 ✚ 1 ↪ 1 ✔ 2 ✎ 1
```

## Verbose mode (binary only)

```bash
__git_status -v   # prints the full JSON state object
# or
git-status-darwin-arm64 -v | jq
```
