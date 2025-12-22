package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// 监听配置文件发生变化
func main() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// 填写要监听的目录或者文件
	path := "E:/Workspace/TOOLS/golang-travel/tour/blog-service/configs/config.yaml"
	_ = watcher.Add(path)
	<-done
}
