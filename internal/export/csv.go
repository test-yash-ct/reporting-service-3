package export

import (
	"encoding/csv"
	"io"
	"strings"
	"unicode/utf8"
)

func WritePatientFinance(w io.Writer, rows [][]string) error {
	cw := csv.NewWriter(w)
	for _, row := range rows {
		out := make([]string, len(row))
		for i, cell := range row {
			out[i] = sanitizeCSVCell(cell)
		}
		if err := cw.Write(out); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func sanitizeCSVCell(cell string) string {
	if cell == "" {
		return cell
	}
	r, _ := utf8.DecodeRuneInString(cell)
	switch r {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + cell
	}
	if strings.HasPrefix(cell, "0x09") || strings.HasPrefix(cell, "0x0D") {
		return "'" + cell
	}
	return cell
}
