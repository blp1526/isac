package server

import "testing"

func TestServerString(t *testing.T) {
	tests := []struct {
		unanonymize bool
		status      string
		zoneName    string
		id          string
		name        string
		want        string
	}{
		{
			unanonymize: true,
			id:          "1129XXXXXXXX",
			name:        "foo",
			status:      "down",
			zoneName:    "is1a",
			want:        "is1a 1129XXXXXXXX   down foo",
		},
		{
			unanonymize: false,
			id:          "1129XXXXXXXX",
			name:        "foo",
			status:      "down",
			zoneName:    "is1a",
			want:        "is1a ************   down foo",
		},
	}

	for _, tt := range tests {
		s := &Server{}
		s.ID = tt.id
		s.Name = tt.name
		s.Instance.Status = tt.status
		s.Zone.Name = tt.zoneName

		got := s.String(tt.unanonymize)
		if got != tt.want {
			t.Errorf("tt.unanonymize: %v, got: %v, tt.want: %v", tt.unanonymize, got, tt.want)
		}
	}

}
