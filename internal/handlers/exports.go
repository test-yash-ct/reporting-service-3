package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/middleware"
)

type ExportAPI struct{}

func (e *ExportAPI) Register(r *gin.RouterGroup) {
	r.GET("/exports/operational", e.Operational)
}

func (e *ExportAPI) Operational(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok || !claims.HasRole("export:ops") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	payload := gin.H{
		"appointments": []gin.H{
			{"id": "a1", "status": "scheduled"},
			{"id": "a2", "status": "scheduled"},
		},
		"billing_events": []gin.H{
			{"id": "b1", "amount_cents": 250000},
			{"id": "b2", "amount_cents": 1200},
		},
	}
	c.JSON(http.StatusOK, payload)
}
