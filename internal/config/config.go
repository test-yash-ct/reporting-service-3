package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ListenAddr     string
	ReportRoot     string
	JWTSecret      string
	ReportKey      []byte
	CORSOrigins    []string
	MaxTokenTTLSec int64
}

func Load() (Config, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET must be set")
	}
	key := os.Getenv("REPORT_ENCRYPTION_KEY")
	if len(key) != 32 {
		return Config{}, fmt.Errorf("REPORT_ENCRYPTION_KEY must be exactly 32 bytes")
	}
	return Config{
		ListenAddr:     getenv("LISTEN_ADDR", "127.0.0.1:8082"),
		ReportRoot:     getenv("REPORT_ROOT", "./data/reports"),
		JWTSecret:      secret,
		ReportKey:      []byte(key),
		CORSOrigins:    splitCSV(getenv("CORS_ORIGINS", "")),
		MaxTokenTTLSec: 900,
	}, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
