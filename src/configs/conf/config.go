package conf

import (
	"IM-Service/src/configs/log"
	"github.com/don764372409/go-im-sdk/client"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var (
	DbPath = "../db/wallet.db"
)
var (
	once sync.Once
	Base *BaseConfig
	Conf *Config
)

type BaseConfig struct {
	BaseDir    string
	ApiHost    string
	WsHost     string
	LogSwitch  string
	LogFile    string
	DeviceType string
}

// InitConfig 初始化方法
func InitConfig(baseConfig *BaseConfig) string {
	//初始化日志
	once.Do(func() {
		Base = &BaseConfig{
			LogFile:    baseConfig.BaseDir + "/logs",
			ApiHost:    baseConfig.ApiHost,
			WsHost:     baseConfig.WsHost,
			DeviceType: baseConfig.DeviceType,
		}
		log.InitLog(baseConfig.BaseDir+"/logs", baseConfig.LogSwitch)
		//启动长连接
		c := client.New(Base.WsHost)
		go c.Startup()
		//读取yaml
		configBytes, _ := os.ReadFile("../configs/config.yaml")
		Conf = &Config{}
		_ = yaml.Unmarshal(configBytes, Conf)
	})
	return baseConfig.BaseDir
}

type Config struct {
	Debug  bool
	Logger *Logger
	Data   *Data
}

type Database struct {
	Driver string
	Source string
}

type Data struct {
	Database   *Database
	ExcludeUri []string
	Prime      string
}

type Logger struct {
	Level string
}
