package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/audit"
	"github.com/healthops/reporting-service/internal/auth"
	"github.com/healthops/reporting-service/internal/data"
	"github.com/healthops/reporting-service/internal/export"
	"github.com/healthops/reporting-service/internal/middleware"
	"github.com/healthops/reporting-service/internal/reports"
)

type ReportAPI struct {
	Files *reports.FileStore
	Data  *data.OperationalStore
	Audit *audit.Logger
}

func (a *ReportAPI) Register(r *gin.RouterGroup) {
	r.GET("/reports/:name.csv", a.CSV)
	r.GET("/reports/file", a.File)
	r.POST("/reports/render", a.Render)
}

func (a *ReportAPI) CSV(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok || !claims.HasRole("report:read") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	rows, err := a.Data.FinanceRows(claims.Tenant)
	if err != nil {
		a.emit(c, claims, "report.csv", "load_failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "report_unavailable"})
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=report.csv")
	c.Status(http.StatusOK)
	_ = export.WritePatientFinance(c.Writer, rows)
	a.emit(c, claims, "report.csv", "ok")
}

func (a *ReportAPI) File(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok || !claims.HasRole("report:read") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	path := c.Query("path")
	b, err := a.Files.ReadFile(claims.Tenant, path)
	if err != nil {
		switch {
		case errors.Is(err, reports.ErrPathTraversal), errors.Is(err, reports.ErrInvalidPath):
			a.emit(c, claims, path, "invalid_path")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_path"})
		case errors.Is(err, reports.ErrFileTooLarge):
			a.emit(c, claims, path, "file_too_large")
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file_too_large"})
		case errors.Is(err, reports.ErrDecryptFailed):
			a.emit(c, claims, path, "decrypt_failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "integrity_failed"})
		default:
			a.emit(c, claims, path, "not_found")
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		}
		return
	}
	a.emit(c, claims, path, "ok")
	c.Data(http.StatusOK, "application/octet-stream", b)
}

func (a *ReportAPI) Render(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok || !claims.HasRole("report:read") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req reports.RenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}
	out, err := reports.RenderSummary(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "render_failed"})
		return
	}
	a.emit(c, claims, req.TemplateName, "ok")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, out)
}

func (a *ReportAPI) emit(c *gin.Context, claims auth.Claims, object, outcome string) {
	if a.Audit == nil {
		return
	}
	rid, _ := c.Get("request_id")
	ridStr, _ := rid.(string)
	a.Audit.Emit(audit.Event{
		Actor: claims.Sub, Tenant: claims.Tenant, Action: "report.access",
		ObjectID: object, Outcome: outcome, RequestID: ridStr,
	})
}
