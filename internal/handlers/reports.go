package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/export"
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
	tenant := c.Query("tenant")
	rows := [][]string{
		{"patient_id", "balance_note"},
		{"P-1001", "=1+1"},
		{"P-1002", "+12025550123"},
	}
	c.Header("Content-Type", "text/csv")
	c.Status(http.StatusOK)
	_ = tenant
	_ = export.WritePatientFinance(c.Writer, rows)
}

func (a *ReportAPI) File(c *gin.Context) {
	tenant := c.Query("tenant")
	path := c.Query("path")
	b, err := a.Files.ReadFile(tenant, path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "read_failed"})
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", b)
}

func (a *ReportAPI) Render(c *gin.Context) {
	var req reports.RenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}
	out, err := reports.RenderSummary(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "render_failed"})
		return
	}
	c.String(http.StatusOK, out)
}
