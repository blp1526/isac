package state

type State struct {
	Current string
}

func New() *State {
	s := &State{}
	return s
}

func (s *State) Toggle(given string) {
	if s.Current == given {
		s.Current = ""
		return
	}

	s.Current = given
}
