package input

import "strings"

func PathState(state map[string]bool, path string) bool {
	if len(state) == 0 {
		return false
	}
	return state[path]
}

func SetPathState(state map[string]bool, path string, on bool) (map[string]bool, bool) {
	p := strings.TrimSpace(path)
	if p == "" {
		return state, false
	}
	if state == nil {
		state = make(map[string]bool)
	}
	current := state[p]
	if current == on {
		return state, false
	}
	if on {
		state[p] = true
	} else {
		delete(state, p)
	}
	return state, true
}



