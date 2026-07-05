// backend/internal/scraper/apifox_models.go
package scraper

// ApifoxTreeNode Apifox shared-doc 侧边栏节点
type ApifoxTreeNode struct {
	NodeType string           `json:"nodeType"` // folder / http_request / doc / api_scenario
	Name     string           `json:"name"`
	API      ApifoxAPIInfo    `json:"api,omitempty"`
	Doc      ApifoxDocInfo    `json:"doc,omitempty"`
	Children []ApifoxTreeNode `json:"children,omitempty"`
}

// ApifoxAPIInfo 接口详情
type ApifoxAPIInfo struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	// 详情需另调接口
}

// ApifoxDocInfo 文档详情
type ApifoxDocInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ApifoxSharedDocRoot shared-doc 根响应
type ApifoxSharedDocRoot struct {
	Node ApifoxTreeNode `json:"node"`
}
