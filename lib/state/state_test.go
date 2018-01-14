package state

import "testing"

func TestStateToggle(t *testing.T) {
	tests := []struct {
		current string
		given   string
		want    string
	}{
		{
			current: "",
			given:   "detail",
			want:    "detail",
		},
		{
			current: "detail",
			given:   "detail",
			want:    "",
		},
	}

	for _, tt := range tests {
		s := &State{Current: tt.current}
		s.Toggle(tt.given)

		if s.Current != tt.want {
			t.Errorf("s.Current: %v, tt.want: %v", s.Current, tt.want)
		}
	}
}
