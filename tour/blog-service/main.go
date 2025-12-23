package main

import (
	_ "blog-service/docs"
	"blog-service/global"
	"blog-service/internal/model"
	"blog-service/internal/routers"
	"blog-service/pkg/logger"
	"blog-service/pkg/setting"
	"blog-service/pkg/tracer"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-programming-tour-book/blog-service/pkg/validator"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	port         string
	runMode      string
	config       string
	isVersion    bool
	buildTime    string
	buildVersion string
	gitCommitID  string
)

// 在main之前自动执行 全局变量初始化 =》init 方法 =》main 方法
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}

// @title 博客系统
// @version 1.0
// @description Go 语言编程之旅：一起用 Go 做项目
// @termsOfService https://github.com/go-programming-tour-book
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token
func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouters()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAn")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited gracefully")
}

func setupSetting() error {
	// NewSetting 创建一个 Setting 实例
	setting, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	// 读取配置文件的各个部分到全局变量中
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	// 读取 App 部分的配置
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	// 读取 Database 部分的配置
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	// 读取 JWT 部分的配置
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	// 读取 Email 部分的配置
	err = setting.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}
	global.AppSetting.DefaultContextTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   500,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupValidator() error {
	global.Validator = validator.NewCustomValidator()
	global.Validator.Engine()
	binding.Validator = global.Validator

	return nil
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
