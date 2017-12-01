package row

type Row struct {
	header        int
	separator     int
	movableTop    int
	current       int
	movableBottom int
}

func New() *Row {
	r := &Row{
		header:     0,
		separator:  1,
		movableTop: 2,
		current:    2,
	}
	return r
}
