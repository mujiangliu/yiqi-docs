// backend/internal/scraper/apifox_browser.go
package scraper

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"gorm.io/gorm"
)

// SeedApifoxSiteViaBrowser JSON API 失败时的备选：用无头浏览器渲染抓取
func SeedApifoxSiteViaBrowser(db *gorm.DB, siteID uint) error {
	// 简化实现：仅写入一个说明页，提示需手动补充
	_, err := InsertPage(db, siteID, nil, "intro", "介绍",
		"# 一起粗去 API 说明\n\n"+
			"> 抓取脚本 JSON API 失败，浏览器渲染抓取未完整实现。\n"+
			"> 请在后台手动编辑本页面内容，或访问 https://kuaipao.apifox.cn/ 获取。\n", 0)
	if err != nil {
		return fmt.Errorf("insert fallback intro: %w", err)
	}
	_ = chromedp.NewContext // 防 unused
	return nil
}
