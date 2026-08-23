package data

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const MaxFileSize = 50 << 20

type OperationalStore struct {
	Root string
}

type OperationalExport struct {
	Appointments  []map[string]any `json:"appointments"`
	BillingEvents []map[string]any `json:"billing_events"`
}

func (s *OperationalStore) tenantPath(tenant string) (string, error) {
	if tenant == "" || strings.Contains(tenant, "..") || strings.ContainsAny(tenant, `/\`) {
		return "", errors.New("invalid tenant")
	}
	return filepath.Join(filepath.Clean(s.Root), tenant, "operational.json"), nil
}

func (s *OperationalStore) Load(tenant string) (OperationalExport, error) {
	p, err := s.tenantPath(tenant)
	if err != nil {
		return OperationalExport{}, err
	}
	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return OperationalExport{}, nil
		}
		return OperationalExport{}, err
	}
	if info.Size() > MaxFileSize {
		return OperationalExport{}, errors.New("file too large")
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		return OperationalExport{}, err
	}
	var out OperationalExport
	if err := json.Unmarshal(raw, &out); err != nil {
		return OperationalExport{}, err
	}
	return out, nil
}

func (s *OperationalStore) FinanceRows(tenant string) ([][]string, error) {
	exp, err := s.Load(tenant)
	if err != nil {
		return nil, err
	}
	rows := [][]string{{"patient_id", "balance_note"}}
	for _, ev := range exp.BillingEvents {
		pid, _ := ev["patient_id"].(string)
		note, _ := ev["status"].(string)
		if pid == "" {
			continue
		}
		rows = append(rows, []string{pid, note})
	}
	if len(rows) == 1 {
		rows = append(rows, []string{"none", "no_data"})
	}
	return rows, nil
}
