package main

import (
	"log"
	"time"
)

// func main() {
// 	s := "string_val"
// 	s = strings.Replace(s, "_", " ", -1) // 将_替换成空格字符
// 	s = strings.Title(s)                 // 首字母大写
// 	s = strings.Replace(s, " ", "", -1)  // 将空格替换为空
// 	fmt.Println(s)
// 	// n = -1 (n<0 替换的次数没有限制 若为2只替换前两个)

// 	a := rune(s[0])
// 	fmt.Println(string(unicode.ToLower(a)) + s[1:])

// 	var output []rune
// 	for i, v := range s {
// 		if i == 0 {
// 			fmt.Println(output)
// 			output = append(output, unicode.ToLower(v))
// 			fmt.Println(output)
// 		}

// 	}
// }

func main() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	inputTime := "2029-09-04 12:02:33"
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, inputTime)
	dateTime := time.Unix(t.Unix(), 0).In(location).Format(layout)

	log.Printf("输入时间：%s，输出时间：%s", inputTime, dateTime)
	locations, _ := time.LoadLocation("Asia/Shanghai")
	log.Println(time.Now().In(locations))

	// 2025/12/10 17:21:12 输入时间：2029-09-04 12:02:33，输出时间：2029-09-04 20:02:33
	// 2025/12/10 17:21:12 2025-12-10 17:21:12.5196418 +0800 CST
}
