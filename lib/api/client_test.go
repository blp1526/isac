package api

import "testing"

func TestClienturl(t *testing.T) {
	tests := []struct {
		zone     string
		resource string
		id       string
		want     string
	}{
		{
			zone:     "tk1a",
			resource: "server",
			id:       "",
			want:     "https://secure.sakura.ad.jp/cloud/zone/tk1a/api/cloud/1.1/server",
		},
	}

	for _, tt := range tests {
		client := &Client{}
		got := client.url(tt.zone, tt.resource, tt.id)

		if got != tt.want {
			t.Errorf("got: %v, tt.want: %v", got, tt.want)
		}
	}
}
