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

func (r *Row) Headers() (headers []string) {
	id := fmt.Sprintf("%-12v", "ID")

	headers = []string{
		fmt.Sprintf("Zone %s Status Name", id),
	}

	return headers
}

func (r *Row) Separator() (separator string) {
	separator = "============================="
	return separator
}
