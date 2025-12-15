package cmd

import (
	"log"
	"strings"

	"github.com/golang-travel/tour/internal/word"
	"github.com/spf13/cobra"
)

const (
	ModeUpper                      = iota + 1 // 全部转大写
	ModeLower                                 // 全部转小写
	ModeUnderscoreToUpperCamelCase            // 下划线转大驼峰
	ModeUnderscoreToLowerCamelCase            // 下划线转小驼峰
	ModeCamelCaseToUnderscore                 // 驼峰转下划线
)

var desc = strings.Join([]string{
	"该子命令支持各种单词格式转换，模式如下：",
	"1：全部转大写",
	"2：全部转小写",
	"3：下划线转大写驼峰",
	"4：下划线转小写驼峰",
	"5：驼峰转下划线",
}, "\n")

var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "单词格式转换",
	Long:  desc,
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case ModeUpper:
			content = word.ToUpper(str)
		case ModeLower:
			content = word.ToLower(str)
		case ModeUnderscoreToUpperCamelCase:
			content = word.UnderscoreToUpperCamelCase(str)
		case ModeUnderscoreToLowerCamelCase:
			content = word.UnderscoreToLowerCamelCase(str)
		case ModeCamelCaseToUnderscore:
			content = word.CamelCaseToUnderscore(str)
		default:
			log.Fatalf("暂不支持该转换模式，请执行 help word 查看帮助文档")
		}

		log.Printf("输出结果: %s", content)
	},
}

var str string
var mode int8

func init() {

	// 1:参数需绑定的变量 2:接收该参数的完整命令标志 3:短标识 4:默认值 5:使用说明
	wordCmd.Flags().StringVarP(&str, "str", "s", "", "请输入单词内容")
	wordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "请输入单词转换的模式")
}

// go run cobra/main.go word  -m=1 -s=jiang
// 2025/12/10 15:19:15 输出结果: JIANG
// go run cobra/main.go word  -m=3 -s=jiangChang
// 2025/12/10 15:19:29 输出结果: JiangChang
// go run cobra/main.go word  -m=3 -s=jiangChang
// 2025/12/10 15:19:35 输出结果: JiangChang
// go run cobra/main.go word  -m=4 -s=jiangChang
// 2025/12/10 15:19:38 输出结果: jiangChang
// go run cobra/main.go word  -m=5 -s=jiangChang
// 2025/12/10 15:19:42 输出结果: jiang_chang
