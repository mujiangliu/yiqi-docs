// backend/cmd/server/main.go
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/api"
	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/config"
	"jiaocheng-web/backend/internal/model"
	"jiaocheng-web/backend/internal/store"
	"jiaocheng-web/backend/web"
)

func main() {
	cfg := config.Load()

	// 打开 DB（用 modernc 纯 Go 驱动，但 gorm sqlite 驱动自动选用）
	// 注意：gorm.io/driver/sqlite 默认用 mattn/go-sqlite3（需 CGO）。
	// 为免 CGO，我们用 modernc.org/sqlite 的 DSN：file:...?_pragma=...
	// 但更简单：直接用 gorm 的 sqlite driver，并 import modernc.org/sqlite 以触发 gorm 的纯 Go 模式
	// 实际：gorm.io/driver/sqlite 1.5+ 支持 modernc，DSN 加 "file:" 前缀即可
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Site{}, &model.Page{}, &model.Media{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	st := store.NewStore(db)

	// 初始化 super_admin
	if err := seedAdmin(db, cfg); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	// JWT issuer
	if cfg.JWTSecret == "" {
		log.Fatalf("JWT_SECRET is required")
	}
	issuer := auth.NewJWTIssuer(cfg.JWTSecret)

	// HTTP
	r := gin.Default()
	api.Router(r, st, issuer)
	web.RegisterStatic(r)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run: %v", err)
	}
}

// seedAdmin 启动时确保 super_admin 存在
func seedAdmin(db *gorm.DB, cfg *config.Config) error {
	if cfg.SeedAdminUser == "" || cfg.SeedAdminPass == "" {
		log.Println("warning: SEED_ADMIN_USER/SEED_ADMIN_PASS not set, skip seeding super_admin")
		return nil
	}
	var existing model.User
	err := db.Where("username = ?", cfg.SeedAdminUser).First(&existing).Error
	if err == nil {
		return nil // 已存在
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hash, err := auth.HashPassword(cfg.SeedAdminPass)
	if err != nil {
		return err
	}
	return db.Create(&model.User{
		Username:     cfg.SeedAdminUser,
		PasswordHash: hash,
		Role:         model.RoleSuperAdmin,
	}).Error
}

// 防止 unused
var _ = os.Getenv
