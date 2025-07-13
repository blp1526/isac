package server

import "testing"

func TestServerString(t *testing.T) {
	tests := []struct {
		status   string
		zoneName string
		id       string
		name     string
		want     string
	}{
		{
			id:       "1129XXXXXXXX",
			name:     "foo",
			status:   "down",
			zoneName: "is1a",
			want:     "is1a 1129XXXXXXXX   down foo",
		},
		{
			id:       "2230YYYYYYYY",
			name:     "bar",
			status:   "up",
			zoneName: "tk1a",
			want:     "tk1a 2230YYYYYYYY     up bar",
		},
	}

	for _, tt := range tests {
		s := &Server{}
		s.ID = tt.id
		s.Name = tt.name
		s.Instance.Status = tt.status
		s.Zone.Name = tt.zoneName

		got := s.String()
		if got != tt.want {
			t.Errorf("got: %v, want: %v", got, tt.want)
		}
	}

}
