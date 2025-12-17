package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProxySQLEndpoints []string
	User              string
	Password          string
	DBName            string
	MaxOpenConns      int
	MaxIdleConns      int
	RouterStrategy    string
	HTTPAddr          string
}

func LoadFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error opening config file: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("error parsing yaml config file: %v", err)
	}
	return cfg, nil
}

func Load() Config {
	return Config{
		ProxySQLEndpoints: splitAndTrim(getenv("PROXYSQL_ENDPOINTS", "127.0.0.1:6033")),
		User:              getenv("DB_USER", "app"),
		Password:          getenv("DB_PASSWORD", "app_password"),
		DBName:            getenv("DB_NAME", "appdb"),
		MaxOpenConns:      atoiDefault(getenv("DB_MAX_OPEN_CONNS", "50"), 50),
		MaxIdleConns:      atoiDefault(getenv("DB_MAX_IDLE_CONNS", "25"), 25),
		RouterStrategy:    getenv("ROUTER_STRATEGY", "round_robin"),
		HTTPAddr:          getenv("HTTP_ADDR", ":8080"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func atoiDefault(s string, def int) int {
	var x int
	if _, err := fmt.Sscanf(s, "%d", &x); err != nil {
		return def
	}
	return x
}

func splitAndTrim(s string) []string {
	var out []string
	var part string
	for _, r := range s {
		if r == ',' {
			if part != "" {
				out = append(out, trim(part))
				part = ""
			}
			continue
		}
		part += string(r)
	}
	if part != "" {
		out = append(out, trim(part))
	}
	return out
}

func trim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
