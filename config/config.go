package config

import (
	"fmt"
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

var (
	TargetHost = getEnv("TARGET_HOST", "192.168.0.92")
	TargetPort = getEnvInt("TARGET_PORT", 5010)
	ProxyPort  = getEnvInt("PROXY_PORT", 5000)
	TargetURL  = fmt.Sprintf("http://%s:%d", TargetHost, TargetPort)
)
