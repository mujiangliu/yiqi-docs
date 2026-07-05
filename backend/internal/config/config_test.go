// backend/internal/config/config_test.go
package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DB_PATH")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SEED_ADMIN_USER")
	os.Unsetenv("SEED_ADMIN_PASS")
	os.Unsetenv("SEED_CONTENT_ADMIN_USER")
	os.Unsetenv("SEED_CONTENT_ADMIN_PASS")

	c := Load()
	if c.Port != 8080 {
		t.Errorf("Port = %d, want 8080", c.Port)
	}
	if c.DBPath != "./data.db" {
		t.Errorf("DBPath = %q, want ./data.db", c.DBPath)
	}
	if c.SeedAdminUser != "admin" {
		t.Errorf("SeedAdminUser = %q, want admin", c.SeedAdminUser)
	}
	if c.SeedContentAdminUser != "content" {
		t.Errorf("SeedContentAdminUser = %q, want content", c.SeedContentAdminUser)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DB_PATH", "/tmp/test.db")
	os.Setenv("JWT_SECRET", "secret123")
	os.Setenv("SEED_ADMIN_USER", "root")
	os.Setenv("SEED_ADMIN_PASS", "pass1")
	os.Setenv("SEED_CONTENT_ADMIN_USER", "writer")
	os.Setenv("SEED_CONTENT_ADMIN_PASS", "pass2")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_PATH")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SEED_ADMIN_USER")
		os.Unsetenv("SEED_ADMIN_PASS")
		os.Unsetenv("SEED_CONTENT_ADMIN_USER")
		os.Unsetenv("SEED_CONTENT_ADMIN_PASS")
	}()

	c := Load()
	if c.Port != 9090 {
		t.Errorf("Port = %d, want 9090", c.Port)
	}
	if c.DBPath != "/tmp/test.db" {
		t.Errorf("DBPath = %q, want /tmp/test.db", c.DBPath)
	}
	if c.JWTSecret != "secret123" {
		t.Errorf("JWTSecret = %q, want secret123", c.JWTSecret)
	}
	if c.SeedAdminUser != "root" || c.SeedAdminPass != "pass1" {
		t.Errorf("SeedAdmin = %q/%q, want root/pass1", c.SeedAdminUser, c.SeedAdminPass)
	}
	if c.SeedContentAdminUser != "writer" || c.SeedContentAdminPass != "pass2" {
		t.Errorf("SeedContentAdmin = %q/%q, want writer/pass2", c.SeedContentAdminUser, c.SeedContentAdminPass)
	}
}
