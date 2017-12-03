package row

import "fmt"

type Row struct {
	MovableTop    int
	MovableBottom int
	Current       int
}

func New() *Row {
	r := &Row{}
	return r
}

func (r *Row) Headers(status string, totalServers int) (headers []string) {
	no := fmt.Sprintf("%-3v", "No")
	id := fmt.Sprintf("%-12v", "ID")

	headers = []string{
		fmt.Sprintf("isac Status: %v", status),
		fmt.Sprintf("Total Servers: %v", totalServers),
		fmt.Sprintf(""),
		fmt.Sprintf("%v Zone %v Status Name", no, id),
	}

	return headers
}
