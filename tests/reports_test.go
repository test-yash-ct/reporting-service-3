package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/handlers"
	"github.com/healthops/reporting-service/internal/middleware"
	"github.com/healthops/reporting-service/internal/reports"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestFileDownloadRequiresAuth(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "t1"), 0o700)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := &handlers.ReportAPI{Files: &reports.FileStore{Root: root, Key: []byte("0123456789abcdef0123456789abcdef")}}
	g := r.Group("/v1")
	g.Use(middleware.Authenticate("test-secret", 900))
	api.Register(g)
	req := httptest.NewRequest(http.MethodGet, "/v1/reports/file?path=secret.txt", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestFileDownloadRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	fs := &reports.FileStore{Root: root, Key: []byte("0123456789abcdef0123456789abcdef")}
	if _, err := fs.ResolveTenantPath("t1", "../etc/passwd"); err == nil {
		t.Fatal("expected traversal rejection")
	}
}
