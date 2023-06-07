package main

import (
	"encoding/json"
	"fmt"
	"grender/core/configReader"
	"grender/core/db"
	"grender/core/model"
	"strconv"
)

func main() {
	configReader.InitConfig()
	// 添加任务
	redis := db.RedisUtil{}
	redis.Connect(configReader.Config.Redis)
	for i := 0; i < 2; i++ {
		task := model.Task{}
		page := strconv.Itoa(i)
		url := fmt.Sprintf("https://exercise.kingname.info/exercise_middleware_ip/%s", page)
		task.Url = url
		//task.Xpath = "//body"
		task.TimeOut = 10
		marshal, err := json.Marshal(task)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(marshal))
		redis.Lpush("test", "string(marshal)")
		fmt.Printf("添加：【%s】 \n", url)
	}
	fmt.Println("done")
}
