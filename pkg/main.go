package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func main() {
	state := &gitState{Data: make(map[string]string)}
	if len(os.Args) > 1 && "-v" == os.Args[1] {
		state.Verbose = true
	}
	state.init()
	fmt.Printf("%s", state)
}

type gitState struct {
	Verbose bool

	// Ref data
	Attached bool
	Data     map[string]string
	Hash     string
	Named    bool
	RefName  string
	Tagged   bool
	Upstream bool

	// local state data
	Added     int
	Ahead     int
	Behind    int
	Deleted   int
	Renamed   int
	Staged    int
	Stashed   int
	Total     int
	Unstaged  int
	Untracked int
}

func (state *gitState) init() {
	state.initRefState()
	state.initLocalState()
}

func (state *gitState) String() string {
	origin := "origin"
	position := "master"

	if !state.Upstream {
		origin = "local"
	}
	if "" != state.Hash {
		position = string([]rune(state.Hash)[:10])
	}
	if state.Named {
		position = state.Data["branch"]
	}
	if !state.Attached {
		origin = "detached"
	}
	if state.Tagged {
		origin = "tag"
		position = fmt.Sprintf("%s (%s)", state.Data["tag"], position)
	}

	status := ""
	if state.Untracked > 0 {
		status = fmt.Sprintf("%s …%d", status, state.Untracked)
	}
	if state.Stashed > 0 {
		status = fmt.Sprintf("%s ＊%d", status, state.Stashed)
	}
	if state.Behind > 0 {
		status = fmt.Sprintf("%s ↓%d", status, state.Behind)
	}
	if state.Ahead > 0 {
		status = fmt.Sprintf("%s ↑%d", status, state.Ahead)
	}
	if state.Deleted > 0 {
		status = fmt.Sprintf("%s ✖ %d", status, state.Deleted)
	}
	if state.Added > 0 {
		status = fmt.Sprintf("%s ✚ %d", status, state.Added)
	}
	if state.Renamed > 0 {
		status = fmt.Sprintf("%s ↪ %d", status, state.Renamed)
	}
	if state.Staged > 0 {
		status = fmt.Sprintf("%s ✔ %d", status, state.Staged)
	}
	if state.Unstaged > 0 {
		status = fmt.Sprintf("%s ✎ %d", status, state.Unstaged)
	}

	if state.Verbose {
		tmp, _ := json.MarshalIndent(state, "", "    ")
		fmt.Println(string(tmp))
		return ""
	}

	return fmt.Sprintf("⎇ %s: %s%s", origin, position, status)
	//return fmt.Sprintf(" %s: %s%s", origin, position, status)
}

func (state *gitState) initLocalState() {
	if "" != state.Data["position"] {
		parts := strings.Split(state.Data["position"], "\t")
		state.Ahead, _ = strconv.Atoi(parts[0])
		state.Behind, _ = strconv.Atoi(parts[1])
	}

	// Parse all file-level state from status --porcelain in a single pass.
	// Unstaged count is derived from the Y-column (working tree status),
	// replacing the separate git diff --name-only subprocess call.
	// Staged count only includes files where X=M and Y=' ' (staged
	// modification with no further working-tree changes), fixing the
	// previous formula that could produce negative counts.
	for _, stat := range strings.Split(state.Data["status"], "\n") {
		if "" == stat {
			continue
		}
		state.Total++
		runes := []rune(stat)
		a := string(runes[0])
		b := string(runes[1])
		if "A" == a || "A" == b {
			state.Added++
		}
		if "D" == a || "D" == b {
			state.Deleted++
		}
		if "M" == a && " " == b {
			state.Staged++
		}
		if "R" == a || "R" == b {
			state.Renamed++
		}
		if "?" == a {
			state.Untracked++
		}
		if " " != b && "?" != b {
			state.Unstaged++
		}
	}
}

var loadMux sync.Mutex

func (state *gitState) load(commands map[string][]string) {
	doneCh := make(chan bool)
	for k, cmd := range commands {
		fnK := k
		fnCmd := cmd
		go func() {
			out, err := exec.Command("git", fnCmd...).Output()
			loadMux.Lock()
			if nil == err {
				state.Data[fnK] = strings.TrimRight(string(out), "\t\n' ")
			} else {
				state.Data[fnK] = ""
			}
			loadMux.Unlock()
			doneCh <- true
		}()
	}

	positionFound := false
	for a := 0; a < len(commands)+1; a++ {
		<-doneCh
		if !positionFound {
			// As soon as the hash and upstream data has loaded, lookup
			// the relative position information.
			loadMux.Lock()
			upstream, upOk := state.Data["upstream"]
			hash, hashOk := state.Data["hash"]
			loadMux.Unlock()
			if upOk && hashOk {
				positionFound = true
				go func() {
					cmpRef := "HEAD"
					if "" != upstream {
						cmpRef = upstream
					}
					out, err := exec.Command("git", strings.Split(fmt.Sprintf("rev-list --left-right --count %s...%s", hash, cmpRef), " ")...).Output()
					loadMux.Lock()
					if nil == err {
						state.Data["position"] = strings.Trim(string(out), "\t\n' ")
					} else {
						state.Data["position"] = ""
					}
					loadMux.Unlock()
					doneCh <- true
				}()
			}
		}
	}
}

// refStateCommands are the git commands run in parallel on each invocation.
// Removed:
//   "abbrev": rev-parse --abbrev-ref HEAD       (result was never read)
//   "ref":    rev-parse --symbolic-full-name HEAD (result was never read)
//   "diff":   diff --name-only                   (redundant: Y-column of
//                                                 status --porcelain encodes
//                                                 the same information)
var refStateCommands = map[string][]string{
	"branch":   {"symbolic-ref", "--short", "HEAD"},
	"hash":     {"rev-parse", "HEAD"},
	"tag":      {"describe", "--exact-match", "--tags", "HEAD"},
	"upstream": {"rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"},

	"stash":  {"stash", "list"},
	"status": {"status", "--porcelain"},
}

func (state *gitState) initRefState() {
	state.load(refStateCommands)
	state.Hash = state.Data["hash"]
	if "" != state.Data["upstream"] {
		state.Upstream = true
	}
	if "" != state.Data["branch"] {
		state.Named = true
	}
	if "" != state.Data["tag"] {
		state.Tagged = true
	}
	if "" != state.Data["stash"] {
		state.Stashed = len(strings.Split(state.Data["stash"], "\n"))
	}
	if state.Upstream || state.Named || state.Tagged {
		state.Attached = true
	}
}
