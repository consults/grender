package configReader

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

var Config *Cfg

type MongoCfg struct {
	Url       string `yaml:"url"`
	DB        string `yaml:"db"`
	RenderCol string `yaml:"render"`
}
type RedisCfg struct {
	Url      string `yaml:"url"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	TaskKey  string `yaml:"task_key"`
}
type Proxy struct {
	Open          bool   `yaml:"open"`
	ProxyAddress  string `yaml:"proxy_address"`
	ProxyUser     string `yaml:"proxy_user"`
	ProxyPassword string `yaml:"proxy_password"`
}

type Api struct {
	Port string `yaml:"port"`
}

type Render struct {
	Local      bool   `yaml:"local"`
	PoolSize   int    `yaml:"pool_size"`
	RodAddress string `yaml:"rod_address"`
}

type Cfg struct {
	Render Render
	Proxy  Proxy
	Mongo  MongoCfg
	Redis  RedisCfg
}

func InitConfig() {
	env := os.Getenv("Config")
	if env == "" {
		env = "dev"
	}
	configPath := fmt.Sprintf("config/%s.yml", env)
	configByte, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("找不到配置文件：%s", configPath))
	}
	c := Cfg{}
	Config = &c
	err = yaml.Unmarshal(configByte, Config)
	if err != nil {
		panic("配置文件格式有问题！")
	}
	log.Warning(Config)
}
