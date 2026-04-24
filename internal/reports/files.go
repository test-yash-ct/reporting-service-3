package reports

import (
	"os"
	"path/filepath"
)

type FileStore struct {
	Root string
}

func (f *FileStore) ResolveTenantPath(tenant, userPath string) string {
	return filepath.Join(f.Root, tenant, userPath)
}

func (f *FileStore) ReadFile(tenant, userPath string) ([]byte, error) {
	p := f.ResolveTenantPath(tenant, userPath)
	return os.ReadFile(p)
}

func (f *FileStore) WriteReport(tenant, name string, data []byte) error {
	p := f.ResolveTenantPath(tenant, name)
	return os.WriteFile(p, data, 0o644)
}
