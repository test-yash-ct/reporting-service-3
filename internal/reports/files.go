package reports

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileStore struct {
	Root string
	Key  []byte
}

func (f *FileStore) ResolveTenantPath(tenant, userPath string) (string, error) {
	if tenant == "" || strings.Contains(tenant, "..") || strings.ContainsAny(tenant, `/\`) {
		return "", errors.New("invalid tenant")
	}
	if userPath == "" || filepath.IsAbs(userPath) {
		return "", errors.New("invalid path")
	}
	cleaned := filepath.Clean(userPath)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(os.PathSeparator)) {
		return "", errors.New("path traversal")
	}
	base := filepath.Join(filepath.Clean(f.Root), tenant)
	full := filepath.Join(base, cleaned)
	rel, err := filepath.Rel(base, full)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", errors.New("path traversal")
	}
	return full, nil
}

func (f *FileStore) ReadFile(tenant, userPath string) ([]byte, error) {
	p, err := f.ResolveTenantPath(tenant, userPath)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return f.decrypt(raw)
}

func (f *FileStore) WriteReport(tenant, name string, data []byte) error {
	p, err := f.ResolveTenantPath(tenant, name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	enc, err := f.encrypt(data)
	if err != nil {
		return err
	}
	return os.WriteFile(p, enc, 0o600)
}

func (f *FileStore) encrypt(plain []byte) ([]byte, error) {
	if len(f.Key) != 32 {
		return nil, errors.New("report encryption key required")
	}
	block, err := aes.NewCipher(f.Key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ct := gcm.Seal(nonce, nonce, plain, nil)
	return []byte(base64.RawStdEncoding.EncodeToString(ct)), nil
}

func (f *FileStore) decrypt(raw []byte) ([]byte, error) {
	if len(f.Key) != 32 {
		return nil, errors.New("report encryption key required")
	}
	block, err := aes.NewCipher(f.Key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	bin, err := base64.RawStdEncoding.DecodeString(string(raw))
	if err != nil {
		return nil, errors.New("ciphertext required")
	}
	ns := gcm.NonceSize()
	if len(bin) < ns {
		return nil, errors.New("ciphertext too short")
	}
	return gcm.Open(nil, bin[:ns], bin[ns:], nil)
}
