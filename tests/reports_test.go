package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/handlers"
	"github.com/healthops/reporting-service/internal/obs"
	"github.com/healthops/reporting-service/internal/reports"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(obs.RequestID())
	r.Use(obs.AccessLogger())
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/meta", obs.MetaHandler(obs.Metadata{
		Service:   "reporting-service",
		Version:   "test",
		BuildTime: "2026-01-01T00:00:00Z",
		GitSHA:    "abc123",
	}))
	return r
}

func TestHealth(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestMetaEndpoint(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/meta", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["service"] != "reporting-service" {
		t.Fatalf("service %q", body["service"])
	}
}

func TestRequestIDEcho(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(obs.HeaderRequestID, "report-req-42")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got := w.Header().Get(obs.HeaderRequestID); got != "report-req-42" {
		t.Fatalf("expected echoed request id, got %q", got)
	}
}

func TestFileDownloadUsesPathParam(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "t1"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "t1", "secret.txt"), []byte("x"), 0o644)
	r := setupRouter()
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
