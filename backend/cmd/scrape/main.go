// backend/cmd/scrape/main.go
package main

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/config"
	"jiaocheng-web/backend/internal/model"
	"jiaocheng-web/backend/internal/scraper"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET required")
	}
	if cfg.SeedAdminUser == "" || cfg.SeedAdminPass == "" {
		log.Fatal("SEED_ADMIN_USER/PASS required for scrape")
	}
	if cfg.SeedContentAdminUser == "" || cfg.SeedContentAdminPass == "" {
		log.Fatal("SEED_CONTENT_ADMIN_USER/PASS required for scrape")
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Site{}, &model.Page{}, &model.Media{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// 1. 账号
	superID, contentID, err := scraper.EnsureAccounts(db, scraper.SeedOptions{
		SuperAdminUser:   cfg.SeedAdminUser,
		SuperAdminPass:   cfg.SeedAdminPass,
		ContentAdminUser: cfg.SeedContentAdminUser,
		ContentAdminPass: cfg.SeedContentAdminPass,
	})
	if err != nil {
		log.Fatalf("ensure accounts: %v", err)
	}
	log.Printf("super_admin id=%d, content admin id=%d", superID, contentID)

	// 2. /yiqicuqu-api 站点
	apiSiteID, err := scraper.EnsureSite(db, contentID, "yiqicuqu-api", "一起粗去 API 说明",
		"一起粗去 API 接口文档（抓取自 kuaipao.apifox.cn）", model.SiteStatusPublished)
	if err != nil {
		log.Fatalf("ensure api site: %v", err)
	}
	if err := scraper.SeedApifoxSite(db, apiSiteID); err != nil {
		log.Fatalf("seed api site: %v", err)
	}
	log.Printf("api site seeded, id=%d", apiSiteID)

	// 3. /yiqicuqu-list 站点
	listSiteID, err := scraper.EnsureSite(db, contentID, "yiqicuqu-list", "一起粗去 使用教程",
		"一起粗去小程序/app 使用教程", model.SiteStatusPublished)
	if err != nil {
		log.Fatalf("ensure list site: %v", err)
	}
	if err := scraper.SeedTutorialSite(db, listSiteID); err != nil {
		log.Fatalf("seed list site: %v", err)
	}
	log.Printf("list site seeded, id=%d", listSiteID)

	log.Println("scrape done")
	os.Exit(0)
}
