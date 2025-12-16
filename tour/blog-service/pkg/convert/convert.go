package convert

import "strconv"

// 统一处理接口返回的响应处理方法
// 类型转换
type StrTo string

// 转换为字符串
func (s StrTo) String() string {
	return string(s)
}

// 转换为整数
func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

// 转换为整数，忽略错误
func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

// 转换为32位无符号整数
func (s StrTo) uint32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

// 转换为32位无符号整数，忽略错误
func (s StrTo) MustUint32() uint32 {
	v, _ := s.uint32()
	return v
}
