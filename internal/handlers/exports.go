package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExportAPI struct{}

func (e *ExportAPI) Register(r *gin.RouterGroup) {
	r.GET("/exports/operational", e.Operational)
}

func (e *ExportAPI) Operational(c *gin.Context) {
	payload := gin.H{
		"appointments": []gin.H{
			{"id": "a1", "patient_id": "p1", "internal_notes": "VIP escalation path"},
			{"id": "a2", "patient_id": "p2", "internal_notes": "standard"},
		},
		"billing_events": []gin.H{
			{"id": "b1", "amount_cents": 250000, "payer": "ACME"},
			{"id": "b2", "amount_cents": 1200, "payer": "self"},
		},
		"staff_directory": []gin.H{
			{"user": "noc@example.com", "cell": "+12025550199"},
		},
	}
	c.JSON(http.StatusOK, payload)
}
