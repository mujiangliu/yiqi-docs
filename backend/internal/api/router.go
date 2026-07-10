// backend/internal/api/router.go
package api

import (
	"github.com/gin-gonic/gin"

	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/store"
)

// Router 组装所有路由
func Router(r *gin.Engine, st *store.Store, issuer *auth.JWTIssuer) {
	r.Use(auth.Middleware(issuer))

	api := r.Group("/api")
	{
		// 公开
		pub := NewPublicHandler(st)
		api.GET("/sites/:path", pub.GetSite)
		api.GET("/sites/:path/pages/*pagePath", pub.GetPage)
		api.GET("/media/:hash", pub.GetMedia)

		// 鉴权
		authH := NewAuthHandler(st, issuer)
		api.POST("/auth/login", authH.Login)
		api.POST("/auth/logout", authH.Logout)
		api.GET("/me", authH.Me)

		// 管理端（需登录）
		admin := api.Group("/admin")
		admin.Use(auth.RequireAuth())
		{
			// 站点
			siteH := NewAdminSiteHandler(st)
			admin.GET("/sites", siteH.List)
			admin.POST("/sites", siteH.Create)
			admin.PUT("/sites/:id", siteH.Update)
			admin.DELETE("/sites/:id", siteH.Delete)
			admin.POST("/sites/:id/media", siteH.UploadMedia)

			// 页面
			pageH := NewAdminPageHandler(st)
			admin.GET("/sites/:id/pages", pageH.ListBySite)
			admin.POST("/sites/:id/pages", pageH.Create)
			admin.GET("/pages/:id", pageH.Get)
			admin.PUT("/pages/:id", pageH.Update)
			admin.DELETE("/pages/:id", pageH.Delete)
			admin.POST("/pages/reorder", pageH.Reorder) // 注意：单独路径避免与 :id 冲突
		}

		// 用户管理（仅 super_admin）
		userH := NewAdminUserHandler(st)
		users := api.Group("/admin/users")
		users.Use(auth.RequireRole("super_admin"))
		{
			users.GET("", userH.List)
			users.POST("", userH.Create)
			users.PUT("/:id", userH.Update)
			users.POST("/:id/reset-password", userH.ResetPassword)
			users.DELETE("/:id", userH.Delete)
		}
	}
}
