package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
)

// 接口校验
// 针对入参校验的方法进行了二次封装

// 单个验证错误
type ValidError struct {
	Key     string
	Message string
}

// 错误集合
type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

// 返回逗号分隔的所有错误信息
func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

// 返回错误信息的字符串切片
func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// 核心函数将 HTTP 请求参数绑定到指定的结构体 v
// 绑定成功 返回true nil
// 绑定失败 从上下文获取翻译器 将错误转换为验证错误
func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	return true, nil
}
