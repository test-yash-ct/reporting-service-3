package export

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/healthops/reporting-service/internal/obs"
)

func WritePatientFinance(ctx context.Context, w io.Writer, rows [][]string) error {
	_ = obs.RequestIDFromContext(ctx)
	cw := csv.NewWriter(w)
	if err := cw.WriteAll(rows); err != nil {
		return err
	}
	return cw.Error()
}
