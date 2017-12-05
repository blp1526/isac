package row

import (
	"reflect"
	"testing"
)

func TestRowHeaders(t *testing.T) {
	tests := []struct {
		status       string
		zones        string
		totalServers int
		currentNo    int
		filter       string
		want         []string
	}{
		{
			status:       "OK",
			zones:        "tk1a",
			totalServers: 3,
			currentNo:    1,
			filter:       "git",
			want:         []string{"isac Status: OK", "Selected Zones: tk1a", "Total Servers: 3, Current No.: 1", "Filter: git", "", "Zone ID           Status Name"},
		},
	}

	for _, tt := range tests {
		r := &Row{}
		got := r.Headers(tt.status, tt.zones, tt.totalServers, tt.currentNo, tt.filter)

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %#v, tt.want: %#v", got, tt.want)
		}
	}
}

func TestRowHeaderSize(t *testing.T) {
	tests := []struct {
		want int
	}{
		{
			want: 6,
		},
	}

	for _, tt := range tests {
		r := &Row{}
		got := r.HeadersSize()

		if got != tt.want {
			t.Errorf("got: %v, tt.want: %v", got, tt.want)
		}
	}
}
