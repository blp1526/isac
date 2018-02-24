package row

import (
	"fmt"

	runewidth "github.com/mattn/go-runewidth"
)

// Row shows this TUI tool's row status.
type Row struct {
	MovableBottom int
	Current       int
	CursorX       int
	CursorY       int
}

// New initializes *Row.
func New() *Row {
	r := &Row{}
	return r
}

// Headers returns this TUI tool's headers.
func (r *Row) Headers(message string, zones string, totalServers int, currentNo int, filter string) (headers []string) {
	id := fmt.Sprintf("%-12v", "ID")

	headers = append(headers,
		fmt.Sprintf("isac Message: %v", message),
		fmt.Sprintf("Selected Zones: %v", zones),
		fmt.Sprintf("Total Servers: %v, Current No.: %v", totalServers, currentNo),
		fmt.Sprintf("Filter: %v", filter),
	)

	r.CursorY = len(headers) - 1
	runes := []rune(headers[r.CursorY])
	x := 0
	for _, r := range runes {
		x += runewidth.RuneWidth(r)
	}
	r.CursorX = x

	headers = append(headers,
		fmt.Sprintf(""),
		fmt.Sprintf("Zone %v Status Name", id),
	)
	return headers
}

// HeadersSize shows headers size.
func (r *Row) HeadersSize() (headersSize int) {
	headersSize = len(r.Headers("", "", 0, 0, ""))
	return headersSize
}
