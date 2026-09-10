package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/obs"
	"github.com/healthops/reporting-service/internal/service"
)

type ExportAPI struct {
	Exporter *service.Exporter
}

func (e *ExportAPI) Register(r *gin.RouterGroup) {
	r.GET("/exports/operational", e.Operational)
}

func (e *ExportAPI) Operational(c *gin.Context) {
	requestID := obs.RequestIDFromContext(c.Request.Context())
	tenant := c.Query("tenant")
	if tenant == "" {
		tenant = c.GetHeader(obs.HeaderTenantID)
	}
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
	logExportEvent(requestID, tenant, "operational_export")
	c.JSON(http.StatusOK, payload)
	if e.Exporter != nil {
		e.Exporter.RecordOperationalExport(c.Request.Context(), tenant, "operational")
	}
}

func logExportEvent(requestID, tenant, event string) {
	entry := map[string]string{
		"event":      event,
		"request_id": requestID,
		"tenant":     tenant,
	}
	b, _ := json.Marshal(entry)
	log.New(os.Stdout, "", 0).Println(string(b))
}
