package api

import "testing"

func TestClientUrl(t *testing.T) {
	tests := []struct {
		zone  string
		paths []string
		want  string
	}{
		{
			zone:  "tk1a",
			paths: []string{"server", ""},
			want:  "https://secure.sakura.ad.jp/cloud/zone/tk1a/api/cloud/1.1/server",
		},
	}

	for _, tt := range tests {
		client := &Client{}
		got := client.Url(tt.zone, tt.paths)

		if got != tt.want {
			t.Errorf("got: %v, tt.want: %v", got, tt.want)
		}
	}
}
