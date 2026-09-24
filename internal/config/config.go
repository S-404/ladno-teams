package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string `env:"ENV"`
	HTTPPort      string `env:"HTTP_PORT"`
	DBHost        string `env:"DB_HOST"`
	DBPort        string `env:"DB_PORT"`
	DBName        string `env:"DB_NAME"`
	DBUser        string `env:"DB_USER"`
	DBPass        string `env:"DB_PASS"`
	DBSsl         string `env:"DB_SSL"`
	JwtSecretKey  string `env:"JWT_SECRET_KEY"`
	JwtAccessTTL  int    `env:"JWT_ACCESS_TTL_SECONDS"`
	JwtRefreshTTL int    `env:"JWT_REFRESH_TTL_SECONDS"`
	CookieSecure  bool   `env:"COOKIE_SECURE"`
	CookieMaxAge  int    `env:"COOKIE_MAX_AGE_SECONDS"`
	AdminLogin    string `env:"ADMIN_LOGIN"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
	AdminName     string `env:"ADMIN_NAME"`
	GitReposRoot  string `env:"GIT_REPOS_ROOT"`
}

func LoadEnv() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.Env = getEnv("ENV", "development")
	cfg.HTTPPort = getEnv("HTTP_PORT", "8080")
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("DB_PORT", "5432")
	cfg.DBName = getEnv("DB_NAME", "ladno_teams")
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPass = getEnv("DB_PASS", "password")
	cfg.DBSsl = getEnv("DB_SSL", "disable")
	cfg.JwtSecretKey = getEnv("JWT_SECRET_KEY", "change-me-secret")
	cfg.JwtAccessTTL = getEnvInt("JWT_ACCESS_TTL_SECONDS", 900)
	cfg.JwtRefreshTTL = getEnvInt("JWT_REFRESH_TTL_SECONDS", 86400)
	cfg.CookieSecure = getEnvBool("COOKIE_SECURE", false)
	cfg.CookieMaxAge = getEnvInt("COOKIE_MAX_AGE_SECONDS", 86400)
	cfg.AdminLogin = getEnv("ADMIN_LOGIN", "admin")
	cfg.AdminPassword = getEnv("ADMIN_PASSWORD", "admin")
	cfg.AdminName = getEnv("ADMIN_NAME", "Administrator")
	cfg.GitReposRoot = getEnv("GIT_REPOS_ROOT", "./data/git-repos")

	if err := cfg.ValidateRequired(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) ValidateRequired() error {
	required := map[string]string{
		"DB_HOST":        c.DBHost,
		"DB_USER":        c.DBUser,
		"DB_PASS":        c.DBPass,
		"JWT_SECRET_KEY": c.JwtSecretKey,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("environment variable %s is required", key)
		}
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	switch strings.ToLower(strValue) {
	case "true", "1":
		return true
	case "false", "0":
		return false
	default:
		return defaultValue
	}
}
