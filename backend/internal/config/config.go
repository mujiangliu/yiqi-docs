// backend/internal/config/config.go
package config

import (
	"os"
	"strconv"
)

// Config 持有运行期配置
type Config struct {
	Port                 int
	DBPath               string
	JWTSecret            string
	SeedAdminUser        string
	SeedAdminPass        string
	SeedContentAdminUser string
	SeedContentAdminPass string
}

// Load 从环境变量加载配置，提供默认值
func Load() *Config {
	port, _ := strconv.Atoi(getenvDefault("PORT", "8080"))
	return &Config{
		Port:                 port,
		DBPath:               getenvDefault("DB_PATH", "./data.db"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		SeedAdminUser:        getenvDefault("SEED_ADMIN_USER", "admin"),
		SeedAdminPass:        os.Getenv("SEED_ADMIN_PASS"),
		SeedContentAdminUser: getenvDefault("SEED_CONTENT_ADMIN_USER", "content"),
		SeedContentAdminPass: os.Getenv("SEED_CONTENT_ADMIN_PASS"),
	}
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
