package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Article struct{}

func NewArticle() Article { return Article{} }

func (a Article) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "获取文章"})
}


func (a Article) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "获取文章列表"})
}


func (a Article) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "创建文章"})
}


func (a Article) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "更新文章"})
}


func (a Article) Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "删除文章"})
}
