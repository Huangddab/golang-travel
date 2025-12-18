package routers

import (
	"blog-service/global"
	"blog-service/internal/middleware"
	"blog-service/internal/routers/api"
	v1 "blog-service/internal/routers/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouters() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Translations())
	// Serve swagger UI. Use the default handler so the UI loads the doc.json
	// from the same host and port as the running server (avoids hardcoded port).
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	article := v1.NewArticle()
	tag := v1.NewTag()
	// 文件上传
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	r.StaticFS("/static", http.Dir(global.AppSetting.Upload.UploadSavePath))

	// 登录认证
	r.POST("/auth", api.GetAuth)
	apiv1 := r.Group("/api/v1")
	{
		apiv1.Use(middleware.JWT())
		// 标签管理
		tags := apiv1.Group("/tags")
		{
			tags.POST("", tag.Create)
			tags.DELETE("/:id", tag.Delete)
			tags.PUT("/:id", tag.Update)
			tags.PATCH("/:id/state", tag.Update)
			tags.GET("", tag.List)
		}

		// 文章管理
		articles := apiv1.Group("/articles")
		{
			articles.POST("", article.Create)
			articles.DELETE("/:id", article.Delete)
			articles.PUT("/:id", article.Update)
			articles.PATCH("/:id/state", article.Update)
			articles.GET("/:id", article.Get)
			articles.GET("", article.List)
		}
	}
	return r
}
