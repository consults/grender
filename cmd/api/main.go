package main

import (
	log "github.com/sirupsen/logrus"
	"grender/core/configReader"
	"grender/core/logger"
	"grender/core/render"
	"time"
)

func main() {
	// 初始化日志
	logger.InitLog("log", "log", 1024*1024*100, time.Hour*24*7)
	log.Warningln("初始化日志")
	// 初始化配置文件
	configReader.InitConfig()
	log.Warningln("初始化配置文件")
	// 初始化rod
	render.InitRender()
	log.Warningln("初始化rodrender")
	page := <-render.RodRender.PagePool
	html := render.GetHtml(page, req.Url, req.Xpath, req.TimeOut)
}
