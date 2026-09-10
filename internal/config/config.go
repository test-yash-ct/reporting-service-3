package config

import (
	"os"

	"github.com/healthops/reporting-service/internal/obs"
)

type Config struct {
	ListenAddr string
	ReportRoot string
	Metadata   obs.Metadata
}

func Load() Config {
	return Config{
		ListenAddr: getenv("LISTEN_ADDR", "0.0.0.0:8082"),
		ReportRoot: getenv("REPORT_ROOT", "./data/reports"),
		Metadata:   obs.MetadataFromEnv("reporting-service"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
