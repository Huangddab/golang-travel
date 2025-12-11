package word

import (
	"strings"
	"unicode"
)

// 全部装大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// 全部转小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// 下划线转大驼峰
func UnderscoreToUpperCamelCase(s string) string {
	// a := "string_val"
	s = strings.Replace(s, "_", " ", -1) // 将_替换成空格字符
	s = strings.Title(s)                 // 首字母大写
	s = strings.Replace(s, " ", "", -1)  // 将空格替换为空

	// n = -1 (n<0 替换的次数没有限制 若为2只替换前两个)
	// fmt.Println(s)
	// StringVal
	return s
}

// 下划线转小驼峰
func UnderscoreToLowerCamelCase(s string) string {
	s = UnderscoreToUpperCamelCase(s)
	a := rune(s[0])
	return string(unicode.ToLower(a)) + s[1:]
	// stringVal
}

// 驼峰转下划线
func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, v := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(v))
			continue
		}
		if unicode.IsUpper(v) {
			// 如果v是大写 其前面加_
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(v))
	}

	return string(output)
}
