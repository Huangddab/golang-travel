package errcode

// 公共错误码
var (
	Success                   = NewError(0, "成功")
	ServerError               = NewError(10000000, "服务内部错误")
	InvalidParams             = NewError(10000001, "入参错误")
	NotFound                  = NewError(10000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(10000003, "鉴权失败，找不到对应的AppKey")
	UnauthorizedTokenError    = NewError(10000004, "鉴权失败，Token错误")
	UnauthorizedTokenTimeout  = NewError(10000005, "鉴权失败，Token超时")
	UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败，Token生成失败")
	TooManyRequests           = NewError(10000007, "请求过多，请稍后再试")
)

// 标签模块
var (
	ErrorGetTagListFail   = NewError(20010001, "获取标签列表失败")
	ErrorCreateTagFail    = NewError(20010002, "创建标签失败")
	ErrorUpdateTagFail    = NewError(20010003, "更新标签失败")
	ErrorDeleteTagFail    = NewError(20010004, "删除标签失败")
	ErrorCountTagFail     = NewError(20010005, "统计标签失败")
	ErrorTagAlreadyExists = NewError(20010006, "标签已存在")
)

// 上传模块
var (
	ErrorUploadFileFail = NewError(20030001, "上传文件失败")
)
