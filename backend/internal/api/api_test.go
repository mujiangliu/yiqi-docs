// backend/internal/api/api_test.go
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/model"
	"jiaocheng-web/backend/internal/store"
)

func setupTestApp(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&model.User{}, &model.Site{}, &model.Page{}, &model.Media{})
	assert.NoError(t, err)
	st := store.NewStore(db)
	issuer := auth.NewJWTIssuer("test-secret")

	// 创建三个用户
	hash, _ := auth.HashPassword("pw")
	users := []model.User{
		{Username: "super", PasswordHash: hash, Role: model.RoleSuperAdmin},
		{Username: "a", PasswordHash: hash, Role: model.RoleAdmin},
		{Username: "b", PasswordHash: hash, Role: model.RoleAdmin},
	}
	for i := range users {
		db.Create(&users[i])
	}

	r := gin.New()
	Router(r, st, issuer)
	return r, db
}

func loginAs(t *testing.T, r *gin.Engine, username string) *http.Cookie {
	body, _ := json.Marshal(map[string]string{"username": username, "password": "pw"})
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies)
	return cookies[0]
}

func doWithCookie(t *testing.T, r *gin.Engine, method, path string, body interface{}, cookie *http.Cookie) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAdminIsolation_EndToEnd(t *testing.T) {
	r, _ := setupTestApp(t)
	cookieA := loginAs(t, r, "a")
	cookieB := loginAs(t, r, "b")
	cookieSuper := loginAs(t, r, "super")

	// adminA 创建站点
	w := doWithCookie(t, r, "POST", "/api/admin/sites", map[string]interface{}{
		"path": "site-a", "title": "Site A", "status": "published",
	}, cookieA)
	assert.Equal(t, 201, w.Code)

	// adminB 看不到 adminA 的站点
	w = doWithCookie(t, r, "GET", "/api/admin/sites", nil, cookieB)
	assert.Equal(t, 200, w.Code)
	var resp struct {
		Data []model.Site `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Empty(t, resp.Data)

	// adminA 看到自己的 1 个
	w = doWithCookie(t, r, "GET", "/api/admin/sites", nil, cookieA)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 1)
	siteAID := resp.Data[0].ID

	// adminB 试图改 adminA 的站点 → 404
	w = doWithCookie(t, r, "PUT", "/api/admin/sites/"+itoa(siteAID), map[string]interface{}{"title": "hacked"}, cookieB)
	assert.Equal(t, 404, w.Code)

	// super 看到全部
	w = doWithCookie(t, r, "GET", "/api/admin/sites", nil, cookieSuper)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 1) // 只有 adminA 创建了 1 个
}

func TestPublicOnlySeesPublished(t *testing.T) {
	r, _ := setupTestApp(t)
	cookieA := loginAs(t, r, "a")

	// 创建 draft 站点
	doWithCookie(t, r, "POST", "/api/admin/sites", map[string]interface{}{
		"path": "draft", "title": "Draft", "status": "draft",
	}, cookieA)

	// 公开访问 draft → 404
	w := doWithCookie(t, r, "GET", "/api/sites/draft", nil, nil)
	assert.Equal(t, 404, w.Code)

	// 改为 published
	w = doWithCookie(t, r, "GET", "/api/admin/sites", nil, cookieA)
	var resp struct {
		Data []model.Site `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	doWithCookie(t, r, "PUT", "/api/admin/sites/"+itoa(resp.Data[0].ID), map[string]interface{}{"status": "published"}, cookieA)

	// 公开访问 → 200
	w = doWithCookie(t, r, "GET", "/api/sites/draft", nil, nil)
	assert.Equal(t, 200, w.Code)
}

func TestPublicSiteLoadsPageContentLazily(t *testing.T) {
	r, _ := setupTestApp(t)
	cookieA := loginAs(t, r, "a")

	doWithCookie(t, r, "POST", "/api/admin/sites", map[string]interface{}{
		"path": "docs", "title": "Docs", "status": "published",
	}, cookieA)

	w := doWithCookie(t, r, "GET", "/api/admin/sites", nil, cookieA)
	var sitesResp struct {
		Data []model.Site `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &sitesResp)
	siteID := sitesResp.Data[0].ID

	doWithCookie(t, r, "POST", "/api/admin/sites/"+itoa(siteID)+"/pages", map[string]interface{}{
		"slug": "intro", "title": "Intro", "content_md": "# Heavy content",
	}, cookieA)

	w = doWithCookie(t, r, "GET", "/api/sites/docs", nil, nil)
	assert.Equal(t, 200, w.Code)
	assert.NotContains(t, w.Body.String(), "Heavy content")

	w = doWithCookie(t, r, "GET", "/api/sites/docs/pages/intro", nil, nil)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Heavy content")
}

func TestAdminPageListExcludesContent(t *testing.T) {
	r, _ := setupTestApp(t)
	cookieA := loginAs(t, r, "a")

	doWithCookie(t, r, "POST", "/api/admin/sites", map[string]interface{}{
		"path": "docs", "title": "Docs", "status": "published",
	}, cookieA)

	w := doWithCookie(t, r, "GET", "/api/admin/sites", nil, cookieA)
	var sitesResp struct {
		Data []model.Site `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &sitesResp)
	siteID := sitesResp.Data[0].ID

	w = doWithCookie(t, r, "POST", "/api/admin/sites/"+itoa(siteID)+"/pages", map[string]interface{}{
		"slug": "intro", "title": "Intro", "content_md": "# Heavy admin content",
	}, cookieA)
	assert.Equal(t, 201, w.Code)
	var createResp struct {
		Data model.Page `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	pageID := createResp.Data.ID

	w = doWithCookie(t, r, "GET", "/api/admin/sites/"+itoa(siteID)+"/pages", nil, cookieA)
	assert.Equal(t, 200, w.Code)
	assert.NotContains(t, w.Body.String(), "Heavy admin content")
	assert.Contains(t, w.Body.String(), "Intro")

	w = doWithCookie(t, r, "GET", "/api/admin/pages/"+itoa(pageID), nil, cookieA)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Heavy admin content")
}

func TestUserManagement_OnlySuper(t *testing.T) {
	r, _ := setupTestApp(t)
	cookieA := loginAs(t, r, "a")
	cookieSuper := loginAs(t, r, "super")

	// admin 试图访问用户列表 → 403
	w := doWithCookie(t, r, "GET", "/api/admin/users", nil, cookieA)
	assert.Equal(t, 403, w.Code)

	// super 访问 → 200
	w = doWithCookie(t, r, "GET", "/api/admin/users", nil, cookieSuper)
	assert.Equal(t, 200, w.Code)
}

// itoa 简单整数转字符串
func itoa(n uint) string {
	return fmtInt(int(n))
}

func fmtInt(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
