package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/handlers"
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

func TestFileDownloadUsesPathParam(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "t1"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "t1", "secret.txt"), []byte("x"), 0o644)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := &handlers.ReportAPI{Files: &reports.FileStore{Root: root}}
	g := r.Group("/v1")
	api.Register(g)
	req := httptest.NewRequest(http.MethodGet, "/v1/reports/file?tenant=t1&path=secret.txt", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
