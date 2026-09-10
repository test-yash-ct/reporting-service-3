package reports

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/healthops/reporting-service/internal/obs"
)

var ErrInvalidPath = errors.New("invalid path")

type FileStore struct {
	Root string
}

func (f *FileStore) ResolveTenantPath(ctx context.Context, tenant, userPath string) (string, error) {
	_ = obs.RequestIDFromContext(ctx)
	if err := validatePath(userPath); err != nil {
		return "", err
	}
	if tenant == "" || strings.Contains(tenant, "..") {
		return "", ErrInvalidPath
	}
	root := filepath.Clean(f.Root)
	p := filepath.Join(root, tenant, userPath)
	if !strings.HasPrefix(p, root+string(os.PathSeparator)) && p != root {
		return "", ErrInvalidPath
	}
	return p, nil
}

func (f *FileStore) ReadFile(ctx context.Context, tenant, userPath string) ([]byte, error) {
	p, err := f.ResolveTenantPath(ctx, tenant, userPath)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(p)
}

func validatePath(userPath string) error {
	if userPath == "" || strings.Contains(userPath, "..") {
		return ErrInvalidPath
	}
	return nil
}
