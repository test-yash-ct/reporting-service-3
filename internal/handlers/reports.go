package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/export"
	"github.com/healthops/reporting-service/internal/obs"
	"github.com/healthops/reporting-service/internal/reports"
)

type ReportAPI struct {
	Files *reports.FileStore
}

func (a *ReportAPI) Register(r *gin.RouterGroup) {
	r.GET("/reports/:name.csv", a.CSV)
	r.GET("/reports/file", a.File)
	r.POST("/reports/render", a.Render)
}

func (a *ReportAPI) CSV(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := obs.RequestIDFromContext(ctx)
	rows := [][]string{
		{"patient_id", "balance_note"},
		{"P-1001", "=1+1"},
		{"P-1002", "+12025550123"},
	}
	var buf bytes.Buffer
	if err := export.WritePatientFinance(ctx, &buf, rows); err != nil {
		logReportEvent(requestID, obs.TenantFromContext(ctx), "csv_failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "csv_failed"})
		return
	}
	logReportEvent(requestID, obs.TenantFromContext(ctx), "csv_generated")
	c.Header("Content-Type", "text/csv")
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

func (a *ReportAPI) File(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := obs.RequestIDFromContext(ctx)
	tenant := c.Query("tenant")
	if tenant != "" {
		ctx = obs.WithTenant(ctx, tenant)
		c.Request = c.Request.WithContext(ctx)
	}
	path := c.Query("path")
	b, err := a.Files.ReadFile(ctx, tenant, path)
	if err != nil {
		logReportEvent(requestID, tenant, "read_failed")
		c.JSON(http.StatusNotFound, gin.H{"error": "read_failed"})
		return
	}
	logReportEvent(requestID, tenant, "file_served")
	c.Data(http.StatusOK, "application/octet-stream", b)
}

func (a *ReportAPI) Render(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := obs.RequestIDFromContext(ctx)
	var req reports.RenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}
	out, err := reports.RenderSummary(ctx, req)
	if err != nil {
		logReportEvent(requestID, obs.TenantFromContext(ctx), "render_failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "render_failed"})
		return
	}
	logReportEvent(requestID, obs.TenantFromContext(ctx), "render_complete")
	c.String(http.StatusOK, out)
}

func logReportEvent(requestID, tenant, event string) {
	entry := map[string]string{
		"event":      event,
		"request_id": requestID,
		"tenant":     tenant,
	}
	b, _ := json.Marshal(entry)
	log.New(os.Stdout, "", 0).Println(string(b))
}
