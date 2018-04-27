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
	state := &gitState{data: make(map[string]string)}
	if len(os.Args) > 1 && "-v" == os.Args[1] {
		state.verbose = true
	}
	state.init()
	fmt.Printf("\n%s\n", state)
}

type gitState struct {
	verbose bool

	// Ref data
	attached bool
	data     map[string]string
	hash     string
	named    bool
	refName  string
	tagged   bool
	upstream bool

	// local state data
	added     int
	ahead     int
	behind    int
	deleted   int
	renamed   int
	staged    int
	stashed   int
	total     int
	unstaged  int
	untracked int
}

func (state *gitState) init() {
	state.initRefState()
	state.initLocalState()
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
				state.data[fnK] = strings.Trim(string(out), "\t\n' ")
			} else {
				state.data[fnK] = ""
			}
			loadMux.Unlock()
			doneCh <- true
		}()
	}
	for a := 0; a < len(commands); a++ {
		<-doneCh
	}
}

func (state *gitState) String() string {
	origin := "origin"
	position := "master"

	if !state.upstream {
		origin = "local"
	}
	if "" != state.hash {
		position = string([]rune(state.hash)[:10])
	}
	if state.named {
		position = state.data["branch"]
	}
	if !state.attached {
		origin = "detached"
	}
	if state.tagged {
		origin = "tag"
		position = state.data["tag"]
	}

	//fmt.Printf(`
	//	state.behind: %d, state.ahead: %d, state.deleted: %d, state.added: %d, state.renamed: %d, state.staged: %d, state.unstaged: %d, state.untracked: %d, state.total: %d
	//`, state.behind, state.ahead, state.deleted, state.added, state.renamed, state.staged, state.unstaged, state.untracked, state.total)
	status := ""
	if state.behind > 0 {
		status = fmt.Sprintf("%s ↓%d", status, state.behind)
	}
	if state.ahead > 0 {
		status = fmt.Sprintf("%s ↑%d", status, state.ahead)
	}
	if state.stashed > 0 {
		status = fmt.Sprintf("%s ＊%d", status, state.stashed)
	}
	if state.deleted > 0 {
		status = fmt.Sprintf("%s ✖%d", status, state.deleted)
	}
	if state.added > 0 {
		status = fmt.Sprintf("%s ✚%d", status, state.added)
	}
	if state.renamed > 0 {
		status = fmt.Sprintf("%s ↪%d", status, state.renamed)
	}
	if state.staged > 0 {
		status = fmt.Sprintf("%s ✔%d", status, state.staged)
	}
	if state.unstaged > 0 {
		status = fmt.Sprintf("%s ✎%d", status, state.unstaged)
	}
	if state.untracked > 0 {
		status = fmt.Sprintf("%s …%d", status, state.untracked)
	}
	if state.total > 0 {
		status = fmt.Sprintf("%s #%d", status, state.total)
	}

	if state.verbose {
		tmp, _ := json.MarshalIndent(state.data, "", "    ")
		fmt.Println(string(tmp))
		fmt.Printf(`
data:     %v

// Ref data
attached: %v
hash:     %v
named:    %v
refName:  %v
tagged:   %v
upstream: %v

// local state data
added:     %v
ahead:     %v
behind:    %v
deleted:   %v
renamed:   %v
staged:    %v
stashed:   %v
total:     %v
unstaged:  %v
untracked: %v
			`,
			state.data,
			state.attached,
			state.hash,
			state.named,
			state.refName,
			state.tagged,
			state.upstream,
			state.added,
			state.ahead,
			state.behind,
			state.deleted,
			state.renamed,
			state.staged,
			state.stashed,
			state.total,
			state.unstaged,
			state.untracked,
		)
	}

	return fmt.Sprintf("⎇ %s/%s%s", origin, position, status)
}

var localStateCommands = map[string][]string{
	"diff":   {"diff", "--name-only"},
	"stash":  {"stash", "list"},
	"status": {"status", "--porcelain"},
}

func (state *gitState) initLocalState() {
	cmpRef := "HEAD"
	if state.upstream {
		//cmpRef = "origin"
	}
	localStateCommands["rev-list"] = strings.Split(fmt.Sprintf("rev-list --left-right --count %s...%s", state.hash, cmpRef), " ")
	fmt.Println(localStateCommands["rev-list"])
	state.load(localStateCommands)

	if "" != state.data["rev-list"] {
		parts := strings.Split(state.data["rev-list"], "\t")
		state.ahead, _ = strconv.Atoi(parts[0])
		state.behind, _ = strconv.Atoi(parts[1])
	}

	if "" != state.data["stash"] {
		state.stashed = len(strings.Split(state.data["stash"], "\n"))
	}

	status := strings.Split(state.data["status"], "\n")
	if "" != state.data["diff"] {
		state.unstaged = len(strings.Split(state.data["diff"], "\n"))
	}
	state.total = len(status)
	for _, stat := range status {
		runes := []rune(stat)
		a := string(runes[0])
		b := string(runes[1])
		//fmt.Printf("\nStatus: '%s', a: '%s', b: '%s', diff: '%s'\n", stat, a, b, state.data["diff"])
		if "A" == a || "A" == b {
			state.added++
		}
		if "D" == a || "D" == b {
			state.deleted++
		}
		if "M" == a || "M" == b {
			state.staged++
		}
		if "R" == a || "R" == b {
			state.renamed++
		}
		if "?" == a || "?" == b {
			state.untracked++
		}
	}
	state.staged -= state.unstaged
	if state.staged < 0 {
		state.staged = 0
	}
}

var refStateCommands = map[string][]string{
	"abbrev": {"rev-parse", "--abbrev-ref", "HEAD"},
	"branch": {"symbolic-ref", "--short", "HEAD"},
	"hash":   {"rev-parse", "HEAD"},
	"ref":    {"rev-parse", "--symbolic-full-name", "HEAD"},
	"tag":    {"describe", "--exact-match", "--tags", "HEAD"},
	//"upstream": {"for-each-ref", "--format='%(upstream:short)'", "$(git symbolic-ref -q HEAD 2> /dev/null)"},
	"upstream": {"rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"},
}

func (state *gitState) initRefState() {
	state.load(refStateCommands)
	state.hash = state.data["hash"]
	if "" != state.data["upstream"] {
		state.upstream = true
	}
	if "" != state.data["branch"] {
		state.named = true
	}
	if "" != state.data["tag"] {
		state.tagged = true
	}
	if "" != state.data["stash"] {
		state.stashed = len(strings.Split(state.data["stash"], "\n"))
	}
	if state.upstream || state.named || state.tagged {
		state.attached = true
	}
}
