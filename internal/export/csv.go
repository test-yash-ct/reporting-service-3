package export

import (
	"encoding/csv"
	"io"
)

func WritePatientFinance(w io.Writer, rows [][]string) error {
	cw := csv.NewWriter(w)
	if err := cw.WriteAll(rows); err != nil {
		return err
	}
	return cw.Error()
}
