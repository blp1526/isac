package row

import "fmt"

type Row struct {
	MovableBottom int
	Current       int
}

func New() *Row {
	r := &Row{}
	return r
}

func (r *Row) Headers(status string, zones string, totalServers int, currentNo int) (headers []string) {
	id := fmt.Sprintf("%-12v", "ID")

	headers = []string{
		fmt.Sprintf("isac Status: %v", status),
		fmt.Sprintf("Selected Zones: %v", zones),
		fmt.Sprintf("Total Servers: %v, Current No.: %v", totalServers, currentNo),
		fmt.Sprintf(""),
		fmt.Sprintf("Zone %v Status Name", id),
	}

	return headers
}

func (r *Row) HeadersSize() (headersSize int) {
	headersSize = len(r.Headers("", "", 0, 0))
	return headersSize
}
