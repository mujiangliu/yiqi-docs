// backend/internal/store/store_test.go
package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/model"
)

func newTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&model.User{}, &model.Site{}, &model.Page{}, &model.Media{})
	assert.NoError(t, err)
	return db
}

func seedUsers(t *testing.T, db *gorm.DB) (super uint, adminA uint, adminB uint) {
	users := []model.User{
		{Username: "super", Role: model.RoleSuperAdmin, PasswordHash: "x"},
		{Username: "a", Role: model.RoleAdmin, PasswordHash: "x"},
		{Username: "b", Role: model.RoleAdmin, PasswordHash: "x"},
	}
	for i := range users {
		err := db.Create(&users[i]).Error
		assert.NoError(t, err)
	}
	return users[0].ID, users[1].ID, users[2].ID
}

func TestSiteCRUD_AdminOnlySeesOwn(t *testing.T) {
	db := newTestDB(t)
	superID, adminAID, adminBID := seedUsers(t, db)

	storeA := NewStore(db).ForUser(&model.User{ID: adminAID, Role: model.RoleAdmin})
	storeB := NewStore(db).ForUser(&model.User{ID: adminBID, Role: model.RoleAdmin})
	storeSuper := NewStore(db).ForUser(&model.User{ID: superID, Role: model.RoleSuperAdmin})

	// adminA 创建站点
	siteA, err := storeA.CreateSite(model.Site{
		OwnerID: adminAID, Path: "site-a", Title: "A", Status: model.SiteStatusPublished,
	})
	assert.NoError(t, err)
	assert.NotZero(t, siteA.ID)

	// adminB 创建站点
	_, err = storeB.CreateSite(model.Site{
		OwnerID: adminBID, Path: "site-b", Title: "B", Status: model.SiteStatusPublished,
	})
	assert.NoError(t, err)

	// adminA 列表只看到自己的 1 个
	listA, err := storeA.ListSites()
	assert.NoError(t, err)
	assert.Len(t, listA, 1)
	assert.Equal(t, "site-a", listA[0].Path)

	// adminB 列表只看到自己的 1 个
	listB, err := storeB.ListSites()
	assert.NoError(t, err)
	assert.Len(t, listB, 1)
	assert.Equal(t, "site-b", listB[0].Path)

	// super 看到全部 2 个
	listSuper, err := storeSuper.ListSites()
	assert.NoError(t, err)
	assert.Len(t, listSuper, 2)

	// adminA 试图读 adminB 的站点 → ErrNotFound
	_, err = storeA.GetSiteByPath("site-b")
	assert.ErrorIs(t, err, ErrNotFound)

	// adminA 试图按 ID 读 adminB 的站点 → ErrNotFound
	bSite, _ := storeB.GetSiteByPath("site-b")
	_, err = storeA.GetSiteByID(bSite.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSiteCRUD_AdminCannotMutateOthers(t *testing.T) {
	db := newTestDB(t)
	_, adminAID, adminBID := seedUsers(t, db)

	storeA := NewStore(db).ForUser(&model.User{ID: adminAID, Role: model.RoleAdmin})
	storeB := NewStore(db).ForUser(&model.User{ID: adminBID, Role: model.RoleAdmin})

	siteB, _ := storeB.CreateSite(model.Site{
		OwnerID: adminBID, Path: "site-b", Title: "B", Status: model.SiteStatusPublished,
	})

	// adminA 试图改 adminB 的站点标题 → ErrNotFound
	_, err := storeA.UpdateSite(siteB.ID, map[string]interface{}{"title": "hacked"})
	assert.ErrorIs(t, err, ErrNotFound)

	// adminA 试图删 adminB 的站点 → ErrNotFound
	err = storeA.DeleteSite(siteB.ID)
	assert.ErrorIs(t, err, ErrNotFound)

	// 站点标题未被改
	unchanged, _ := storeB.GetSiteByID(siteB.ID)
	assert.Equal(t, "B", unchanged.Title)
}

func TestPathUnique(t *testing.T) {
	db := newTestDB(t)
	_, adminAID, _ := seedUsers(t, db)
	storeA := NewStore(db).ForUser(&model.User{ID: adminAID, Role: model.RoleAdmin})

	_, err := storeA.CreateSite(model.Site{OwnerID: adminAID, Path: "dup", Title: "1"})
	assert.NoError(t, err)

	_, err = storeA.CreateSite(model.Site{OwnerID: adminAID, Path: "dup", Title: "2"})
	assert.Error(t, err) // 唯一约束
}

func TestPageTreeOwnerScoped(t *testing.T) {
	db := newTestDB(t)
	_, adminAID, adminBID := seedUsers(t, db)
	storeA := NewStore(db).ForUser(&model.User{ID: adminAID, Role: model.RoleAdmin})
	storeB := NewStore(db).ForUser(&model.User{ID: adminBID, Role: model.RoleAdmin})

	siteA, _ := storeA.CreateSite(model.Site{OwnerID: adminAID, Path: "a", Title: "A", Status: model.SiteStatusPublished})
	siteB, _ := storeB.CreateSite(model.Site{OwnerID: adminBID, Path: "b", Title: "B", Status: model.SiteStatusPublished})

	// adminA 在自己站点下创建页面
	pageA, err := storeA.CreatePage(model.Page{SiteID: siteA.ID, Slug: "intro", Title: "Intro", ContentMD: "# A"})
	assert.NoError(t, err)

	// adminB 试图在 adminA 的站点下创建页面 → ErrNotFound（站点对 B 不可见）
	_, err = storeB.CreatePage(model.Page{SiteID: siteA.ID, Slug: "x", Title: "x"})
	assert.ErrorIs(t, err, ErrNotFound)

	// adminB 试图读取 adminA 的页面 → ErrNotFound
	_, err = storeB.GetPageByID(pageA.ID)
	assert.ErrorIs(t, err, ErrNotFound)

	// adminB 试图更新 adminA 的页面 → ErrNotFound
	_, err = storeB.UpdatePage(pageA.ID, map[string]interface{}{"title": "hacked"})
	assert.ErrorIs(t, err, ErrNotFound)

	// adminA 列出自己站点的页面树
	pages, err := storeA.ListPagesBySite(siteA.ID)
	assert.NoError(t, err)
	assert.Len(t, pages, 1)
	assert.Equal(t, "# A", pages[0].ContentMD)

	// 列表树接口不加载正文
	tree, err := storeA.ListPageTreeBySite(siteA.ID)
	assert.NoError(t, err)
	assert.Len(t, tree, 1)
	assert.Equal(t, "Intro", tree[0].Title)
	assert.Empty(t, tree[0].ContentMD)
	_ = siteB // 避免 unused
}

func TestDeleteSiteCascades(t *testing.T) {
	db := newTestDB(t)
	_, adminAID, _ := seedUsers(t, db)
	storeA := NewStore(db).ForUser(&model.User{ID: adminAID, Role: model.RoleAdmin})

	site, _ := storeA.CreateSite(model.Site{OwnerID: adminAID, Path: "a", Title: "A", Status: model.SiteStatusPublished})
	_, _ = storeA.CreatePage(model.Page{SiteID: site.ID, Slug: "p1", Title: "P1"})
	_, _ = storeA.CreatePage(model.Page{SiteID: site.ID, Slug: "p2", Title: "P2"})

	err := storeA.DeleteSite(site.ID)
	assert.NoError(t, err)

	pages, _ := storeA.ListPagesBySite(site.ID)
	assert.Empty(t, pages)
}

func TestPublicSiteOnlyPublished(t *testing.T) {
	db := newTestDB(t)
	_, adminAID, _ := seedUsers(t, db)
	storeA := NewStore(db).ForUser(&model.User{ID: adminAID, Role: model.RoleAdmin})

	// draft 站点
	storeA.CreateSite(model.Site{OwnerID: adminAID, Path: "draft", Title: "D", Status: model.SiteStatusDraft})
	// published 站点
	storeA.CreateSite(model.Site{OwnerID: adminAID, Path: "pub", Title: "P", Status: model.SiteStatusPublished})

	// 公开视角：只能看到 published
	publicStore := NewStore(db).Public()
	pub, err := publicStore.GetPublishedSiteByPath("pub")
	assert.NoError(t, err)
	assert.Equal(t, "P", pub.Title)

	_, err = publicStore.GetPublishedSiteByPath("draft")
	assert.ErrorIs(t, err, ErrNotFound)
}
