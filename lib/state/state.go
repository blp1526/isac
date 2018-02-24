package state

// State show this TUI tool's state.
type State struct {
	Current string
}

// New initializes *State.
func New() *State {
	s := &State{}
	return s
}

// Toggle changes current state.
func (s *State) Toggle(given string) {
	if s.Current == given {
		s.Current = ""
		return
	}

	s.Current = given
}
