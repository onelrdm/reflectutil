package xlsxutil

import (
	"github.com/tealeg/xlsx"
)

type CellWriter struct {
	*xlsx.Cell
}

func (r *CellWriter) WriteString(s string) error {
	r.Cell.Value = s
	return nil
}
