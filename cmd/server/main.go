package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/config"
	"github.com/healthops/reporting-service/internal/handlers"
	"github.com/healthops/reporting-service/internal/reports"
)

func main() {
	cfg := config.Load()
	if err := os.MkdirAll(cfg.ReportRoot, 0o755); err != nil {
		log.Fatalf("mkdir report root: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	fs := &reports.FileStore{Root: filepath.Clean(cfg.ReportRoot)}
	v1 := r.Group("/v1")
	(&handlers.ReportAPI{Files: fs}).Register(v1)
	(&handlers.ExportAPI{}).Register(v1)

	srv := &http.Server{Addr: cfg.ListenAddr, Handler: r}
	go func() {
		log.Printf("listening on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
