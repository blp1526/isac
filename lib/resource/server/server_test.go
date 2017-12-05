package server

import "testing"

func TestServerString(t *testing.T) {
	tests := []struct {
		showServerID bool
		status       string
		zoneName     string
		id           string
		name         string
		want         string
	}{
		{
			showServerID: true,
			id:           "1129XXXXXXXX",
			name:         "foo",
			status:       "down",
			zoneName:     "is1a",
			want:         "is1a 1129XXXXXXXX   down foo",
		},
		{
			showServerID: false,
			id:           "1129XXXXXXXX",
			name:         "foo",
			status:       "down",
			zoneName:     "is1a",
			want:         "is1a ************   down foo",
		},
	}

	for _, tt := range tests {
		s := &Server{}
		s.ID = tt.id
		s.Name = tt.name
		s.Instance.Status = tt.status
		s.Zone.Name = tt.zoneName

		got := s.String(tt.showServerID)
		if got != tt.want {
			t.Errorf("tt.showServerID: %v, got: %v, tt.want: %v", tt.showServerID, got, tt.want)
		}
	}

}
