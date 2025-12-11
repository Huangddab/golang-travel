package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
)

// 1.1 flag基本使用
// func main() {
// 	var name string
// 	// 对命令行参数的解析和绑定 参数：name 默认值：GO 语言编程之旅 备注：
// 	flag.StringVar(&name, "name", "GO 语言编程之旅", "帮助消息")
// 	flag.StringVar(&name, "n", "GO 语言编程之旅", "帮助消息")

// 	flag.Parse()

// 	log.Printf("name: %s", name)
// }

// var name string

// 1.3 子命令的实现
// func main() {
// 	flag.Parse() // 将命令行解析为定义的标志

// 	args := flag.Args() // 获取命令行非标志参数

// 	if len(args) <= 0 {
// 		return
// 	}

// 	switch args[0] {
// 	case "go":
// 		goCmd := flag.NewFlagSet("go", flag.ExitOnError)
// 		// NewFlagSet 返回指定名称和错误处理属性的空命令集给我们去使用，相当于创建一个新的命令集了
// 		// 第二个参数 用于指定异常错误的情况处理
// 		// const (
// 		// 	ContinueOnError ErrorHandling = iota // Return a descriptive error.
// 		// 	ExitOnError                          // Call os.Exit(2) or for -h/-help Exit(0).
// 		// 	PanicOnError                         // Call panic with a descriptive error.
// 		// )

// 		goCmd.StringVar(&name, "name", "GO 语言", "帮助信息")
// 		_ = goCmd.Parse(args[1:])
// 	case "php":
// 		goCmd := flag.NewFlagSet("php", flag.ExitOnError)
// 		goCmd.StringVar(&name, "name", "PHP 语言", "帮助信息")
// 		_ = goCmd.Parse(args[1:])
// 	}
// 	fmt.Println("name:", name)
// }

// 1.5 可以自定义参数类型

type Name string

// 实现了Value和Set方法
func (i *Name) String() string {
	return fmt.Sprint(*i)
}

func (i *Name) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("name flag already set")
	}
	*i = Name("huang:" + value)
	return nil
}
func main() {
	var name Name
	flag.Var(&name, "name", "帮助信息")
	flag.Parse()
	log.Printf("name:%s", name)

}

// PS E:\Workspace\TOOLS\golang-travel\tour>  go run main.go --name=golang
// s: --name=golang
// 2025/12/10 13:46:32 name:huang:golang
