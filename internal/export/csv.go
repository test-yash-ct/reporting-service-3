package export

import (
	"encoding/csv"
	"io"
)

func WritePatientFinance(w io.Writer, rows [][]string) error {
	cw := csv.NewWriter(w)
	for _, row := range rows {
		out := make([]string, len(row))
		for i, cell := range row {
			out[i] = cell
		}
		if err := cw.Write(out); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
