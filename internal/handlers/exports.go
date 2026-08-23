package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/audit"
	"github.com/healthops/reporting-service/internal/data"
	"github.com/healthops/reporting-service/internal/middleware"
)

type ExportAPI struct {
	Data  *data.OperationalStore
	Audit *audit.Logger
}

func (e *ExportAPI) Register(r *gin.RouterGroup) {
	r.GET("/exports/operational", e.Operational)
}

func (e *ExportAPI) Operational(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok || !claims.HasRole("export:ops") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	exp, err := e.Data.Load(claims.Tenant)
	if err != nil {
		e.emit(c, claims.Sub, claims.Tenant, "operational", "load_failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export_unavailable"})
		return
	}
	payload := gin.H{
		"tenant":         claims.Tenant,
		"appointments":   exp.Appointments,
		"billing_events": exp.BillingEvents,
	}
	e.emit(c, claims.Sub, claims.Tenant, "operational", "ok")
	c.JSON(http.StatusOK, payload)
}

func (e *ExportAPI) emit(c *gin.Context, actor, tenant, object, outcome string) {
	if e.Audit == nil {
		return
	}
	rid, _ := c.Get("request_id")
	ridStr, _ := rid.(string)
	e.Audit.Emit(audit.Event{
		Actor: actor, Tenant: tenant, Action: "export.access",
		ObjectID: object, Outcome: outcome, RequestID: ridStr,
	})
}
