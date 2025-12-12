package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-travel/tour/blog-service/global"
	"github.com/golang-travel/tour/blog-service/internal/model/routers"
	"github.com/golang-travel/tour/blog-service/pkg/setting"
)

func init() {
	error := setupSetting()
	if error != nil {
		log.Fatalf("init.setupSetting err: %v", error)
	}

}

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
	log.Println("Server starting on :8080")
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}
