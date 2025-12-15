package v1

import "github.com/gin-gonic/gin"

type Tag struct {
}

// 构造函数
// 生产和返回一个初始化好的结构体实例
func NewTag() Tag {
	return Tag{}
}

// 这里使用的是t Tag类型的值接收者 意味着每次调用方法时
// 都会得到Tag类型的一个副本 而不是对原始结构体的引用
// 适用于不需要修改结构体内部状态的方法

// 更常见的是使用指针接收者 以避免复制开销 并允许方法修改结构体的状态
func (t Tag) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "获取标签",
	})
}

func (t Tag) List(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "获取标签列表",
	})
}

func (t Tag) Create(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建标签",
	})
}

func (t Tag) Update(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "获取标签",
	})
}

func (t Tag) Delete(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "获取标签",
	})
}
