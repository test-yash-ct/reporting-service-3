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
	"github.com/healthops/reporting-service/internal/audit"
	"github.com/healthops/reporting-service/internal/config"
	"github.com/healthops/reporting-service/internal/data"
	"github.com/healthops/reporting-service/internal/handlers"
	"github.com/healthops/reporting-service/internal/middleware"
	"github.com/healthops/reporting-service/internal/reports"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	_ = os.MkdirAll(cfg.ReportRoot, 0o700)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.RequestLogger())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/readyz", func(c *gin.Context) {
		if err := reports.CheckStorage(cfg.ReportRoot); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage_unavailable"})
			return
		}
		c.String(http.StatusOK, "ok")
	})

	fs := &reports.FileStore{Root: filepath.Clean(cfg.ReportRoot), Key: cfg.ReportKey}
	ds := &data.OperationalStore{Root: filepath.Clean(cfg.ReportRoot)}
	al := audit.New()
	v1 := r.Group("/v1")
	v1.Use(middleware.Authenticate(cfg.JWTSecret, cfg.MaxTokenTTLSec))
	v1.Use(middleware.RateLimit(60, time.Minute))
	v1.Use(middleware.MaxBodyBytes(1 << 20))
	(&handlers.ReportAPI{Files: fs, Data: ds, Audit: al}).Register(v1)
	(&handlers.ExportAPI{Data: ds, Audit: al}).Register(v1)

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
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
	_ = srv.Shutdown(shutdownCtx)
}
