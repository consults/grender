package main

import (
	"encoding/json"
	"github.com/go-rod/rod"
	log "github.com/sirupsen/logrus"
	"grender/core/configReader"
	"grender/core/db"
	"grender/core/logger"
	"grender/core/model"
	"grender/core/render"
	"time"
)

func getTask(taskQueue chan model.Task, redis *db.RedisUtil) {
	for {
		task := redis.Lpop()
		if task == "" {
			time.Sleep(time.Second * 1)
			continue
		}
		modelTask := model.Task{}
		modelTaskErr := json.Unmarshal([]byte(task), &modelTask)
		if modelTaskErr != nil {
			log.Errorln(modelTaskErr.Error())
			continue
		}

		taskQueue <- modelTask
		log.Infof("添加任务：【%s】\n", modelTask.Url)
	}
}
func renderHtml(page *rod.Page, task model.Task, mongo *db.MongoUtil) {
	defer func() {
		render.RodRender.PagePool <- page
	}()
	if len(task.Cookies) > 0 {
		addCookErr := render.AddCookies(page, task.Cookies)
		if addCookErr != nil {
			log.Errorf("添加cookies失败：【%s】 \n", task.Url)
		}
	}
	renderDone := render.WaitLoadElement(page, task.Url, task.Xpath, task.TimeOut)
	html := ""
	if renderDone {
		html = page.MustHTML()
	}
	mongo.InsertOne(task.Url, html, task.Xpath, renderDone)
	log.Infof("渲染完毕：【%s】 \n", task.Url)
}
func runTask(taskQueue chan model.Task, render *render.Render, mongo *db.MongoUtil) {
	for {
		task := <-taskQueue
		page := <-render.PagePool
		go renderHtml(page, task, mongo)
	}
}

func main() {
	logger.InitLog("log", "log", 1024*1024*100, time.Hour*24*7)
	configReader.InitConfig()
	log.Warningln("初始配置文件")
	render.InitRender()
	log.Warningln("初始化渲染池")
	mongo := db.MongoUtil{}
	mongo.Connect(configReader.Config.Mongo)
	redis := db.RedisUtil{}
	redis.Connect(configReader.Config.Redis)
	TaskQueue := make(chan model.Task, configReader.Config.Render.PoolSize*2)
	go getTask(TaskQueue, &redis)
	go runTask(TaskQueue, render.RodRender, &mongo)
	select {}
}
