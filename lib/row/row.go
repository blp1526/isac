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
	no := fmt.Sprintf("%-3v", "No")
	id := fmt.Sprintf("%-12v", "ID")

	headers = []string{
		fmt.Sprintf("%s Zone %s Status Name", no, id),
	}

	return headers
}

func (r *Row) Separator() (separator string) {
	separator = "============================="
	return separator
}
