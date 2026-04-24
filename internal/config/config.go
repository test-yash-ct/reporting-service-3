package config

import "os"

type Config struct {
	ListenAddr string
	ReportRoot string
}

func Load() Config {
	return Config{
		ListenAddr:  getenv("LISTEN_ADDR", "0.0.0.0:8082"),
		ReportRoot:  getenv("REPORT_ROOT", "./data/reports"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
