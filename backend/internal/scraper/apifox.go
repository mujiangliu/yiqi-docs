// backend/internal/scraper/apifox.go
package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	apifoxSharedDocTokenEnv = "APIFOX_SHARED_DOC_TOKEN"
	apifoxUserAgent         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
)

// SeedApifoxSite 抓取 kuaipao.apifox.cn 写入站点
func SeedApifoxSite(db *gorm.DB, siteID uint) error {
	log.Println("fetching apifox shared-doc tree...")
	root, err := fetchApifoxTree()
	if err != nil {
		log.Printf("apifox JSON API failed: %v, falling back to browser scrape", err)
		return SeedApifoxSiteViaBrowser(db, siteID)
	}

	// 清空旧页面
	if err := ClearPages(db, siteID); err != nil {
		return fmt.Errorf("clear pages: %w", err)
	}

	// 先写入一个介绍页（首页）
	introMD := generateIntroMarkdown()
	_, err = InsertPage(db, siteID, nil, "intro", "介绍", introMD, 0)
	if err != nil {
		return fmt.Errorf("insert intro: %w", err)
	}

	// 遍历树，每个 folder → 一级页面，每个接口/文档 → 子页面
	sortIdx := 1
	for _, child := range root.Children {
		id, err := insertTreeNode(db, siteID, nil, child, sortIdx)
		if err != nil {
			log.Printf("insert node %q: %v", child.Name, err)
		}
		if id > 0 {
			sortIdx++
		}
	}
	return nil
}

func fetchApifoxTree() (*ApifoxTreeNode, error) {
	baseURL, err := apifoxSharedDocBaseURL()
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", baseURL+"/http-requests", nil)
	req.Header.Set("User-Agent", apifoxUserAgent)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("apifox tree status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var node ApifoxTreeNode
	if err := json.Unmarshal(body, &node); err != nil {
		return nil, fmt.Errorf("unmarshal tree: %w", err)
	}
	return &node, nil
}

// insertTreeNode 递归插入树节点
func insertTreeNode(db *gorm.DB, siteID uint, parentID *uint, node ApifoxTreeNode, sortIdx int) (uint, error) {
	if node.Name == "" {
		return 0, nil
	}
	slug := slugify(node.Name)
	switch node.NodeType {
	case "folder", "":
		// folder 作为父页面
		pageID, err := InsertPage(db, siteID, parentID, slug, node.Name, "", sortIdx)
		if err != nil {
			return 0, err
		}
		pid := &pageID
		childSort := 1
		for _, child := range node.Children {
			id, err := insertTreeNode(db, siteID, pid, child, childSort)
			if err != nil {
				log.Printf("insert child %q: %v", child.Name, err)
			}
			if id > 0 {
				childSort++
			}
		}
		return pageID, nil
	case "http_request":
		md, err := fetchAPIDetailMarkdown(node.API)
		if err != nil {
			log.Printf("fetch api detail %q: %v", node.Name, err)
			md = fmt.Sprintf("# %s\n\n（接口详情获取失败：%v）\n", node.Name, err)
		}
		return InsertPage(db, siteID, parentID, slug, node.Name, md, sortIdx)
	case "doc":
		md, err := fetchDocMarkdown(node.Doc.ID)
		if err != nil {
			log.Printf("fetch doc detail %q: %v", node.Name, err)
			md = fmt.Sprintf("# %s\n\n（文档详情获取失败：%v）\n", node.Name, err)
		}
		return InsertPage(db, siteID, parentID, slug, node.Name, md, sortIdx)
	}
	return 0, nil
}

// fetchAPIDetailMarkdown 抓取单个接口详情并转 Markdown
func fetchAPIDetailMarkdown(api ApifoxAPIInfo) (string, error) {
	baseURL, err := apifoxSharedDocBaseURL()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/http-requests/%s", baseURL, api.ID)
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", apifoxUserAgent)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("api detail status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 简化：把原始 JSON 转成 Markdown 代码块（实际 Apifox 返回结构复杂，
	// 这里先把原始内容展示出来，后续可解析为更友好的 Markdown）
	return fmt.Sprintf("# %s\n\n`%s %s`\n\n## 原始响应\n\n```json\n%s\n```\n",
		api.Name, api.Method, api.Path, string(body)), nil
}

// fetchDocMarkdown 抓取文档详情
func fetchDocMarkdown(docID string) (string, error) {
	baseURL, err := apifoxSharedDocBaseURL()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/docs/%s", baseURL, docID)
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", apifoxUserAgent)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("doc detail status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("## 原始文档\n\n```json\n%s\n```\n", string(body)), nil
}

func apifoxSharedDocBaseURL() (string, error) {
	token := strings.TrimSpace(os.Getenv(apifoxSharedDocTokenEnv))
	if token == "" {
		return "", fmt.Errorf("%s is not set", apifoxSharedDocTokenEnv)
	}
	return "https://kuaipao.apifox.cn/api/v1/shared-docs/" + token, nil
}

// generateIntroMarkdown 生成介绍页 Markdown
func generateIntroMarkdown() string {
	return `# 一起粗去 API 说明

本文档抓取自 [kuaipao.apifox.cn](https://kuaipao.apifox.cn/)，由后台抓取脚本生成。

## 快跑 AI 接口聚合管理服务

### 核心优势

- 快跑 API 系统有专业的技术团队维护稳定，自动切换号池，专业的技术客服，选择快跑就是选择稳定持续低价 Token。
- 可开发票个人企业普调专票。
- 公司位于深圳市宝安区福海街道，可联系运营随时欢迎洽谈。
- 国内域名：kuaipao.pro
- 国际域名：kuaipao.ai

### 便捷访问

- 无需科学上网，全球直连
- 无封号风险
- 连接速度是官方的 1200 倍
- 覆盖全球七大地区（美国、日本、韩国、英国、香港、菲律宾、俄罗斯）

### 功能特点

- 完善的模型权限
- 支持最新模型直接使用
- 一个 API key 全模型通用
- API key 可设定使用时间和额度
- 100% 保值换绑使用

### 功能对比表

| 功能 | 快跑 API | 官方 API |
|---|---|---|
| 支持 GPT 等模型 | 支持各类型 | 需要账号有 4.0 权限 |
| 最高调用速度 | 支持 | 需绑卡付费 48 小时后 |
| 多账号高并发开发 | 数百个账号 | 单个账号 API 有限制 |
| OpenAI 账号要求 | 无需注册 | 需科学上网和绑定国外手机 |
| 额度有效期 | 永不过期 | 三个月到期 |
| 风控问题 | 0 封号 | 随时无故封号 |
| 使用记录查看 | 实时查看，保留 30 天 | 只能看到延迟总消耗 |
| 代理访问要求 | 无需代理 | 需要在可支持的地区使用 |
| 计费规则 | 折扣价 | 原价 |
| 接口地址 | https://kuaipao.ai | https://api.openai.com |
| 资源整合 | 完全兼容各平台接口协议 | - |

### 技术优势

- 采用企业高速链
- 无需路由二次保留数据
- 稳定纯净 API 源头
- 支持 OpenAI 接口协议
- 0 开发基础可直接对接各种应用

> 注：本页面内容为抓取脚本生成，可在后台编辑修改。
`
}

// slugify 简单 slug 转换
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "&", "and")
	// 保留中文字符（前端按 path 匹配）
	return s
}
