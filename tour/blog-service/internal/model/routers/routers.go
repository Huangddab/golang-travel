package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/golang-travel/tour/blog-service/internal/model/routers/api/v1"
)

func NewRouters() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	article := v1.NewArticle()
	tag := v1.NewTag()

	apiv1 := r.Group("/api/v1")
	{
		// 标签管理
		tags := apiv1.Group("/tags")
		{
			tags.POST("", tag.Create)
			tags.DELETE("/:id", tag.Delete)
			tags.PUT("/:id", tag.Update)
			tags.PATCH("/:id/state", tag.Update)
			tags.GET("/", tag.List)
		}

		// 文章管理
		articles := apiv1.Group("/articles")
		{
			articles.POST("", article.Create)
			articles.DELETE("/:id", article.Delete)
			articles.PUT("/:id", article.Update)
			articles.PATCH("/:id/state", article.Update)
			articles.GET("/:id", article.Get)
			articles.GET("/", article.List)
		}

	}
	return r
}
