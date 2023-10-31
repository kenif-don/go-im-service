package conf

import (
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"gopkg.in/yaml.v3"
	"im-sdk/client"
	"os"
	"sync"
)

var (
	DbPath = "../db/wallet.db"
)
var (
	once      sync.Once
	Base      *BaseConfig
	Conf      *Config
	LoginInfo *LoginInfoMode
)

type LoginInfoMode struct {
	Token string
	User  *entity.User
}
type BaseConfig struct {
	BaseDir    string
	ApiHost    string
	WsHost     string
	LogSwitch  string
	LogFile    string
	DeviceType string
}

// InitConfig 初始化方法
func InitConfig(baseConfig *BaseConfig) (string, error) {
	//初始化日志
	once.Do(func() {
		Base = &BaseConfig{
			LogFile:    baseConfig.BaseDir + "/logs",
			ApiHost:    baseConfig.ApiHost,
			WsHost:     baseConfig.WsHost,
			DeviceType: baseConfig.DeviceType,
		}
		log.InitLog(baseConfig.BaseDir+"/logs", baseConfig.LogSwitch)
		////启动长连接
		c := client.New(Base.WsHost)
		err := c.Startup()
		if err != nil {
			log.Error(err)
		}
		//读取yaml
		configBytes, _ := os.ReadFile("../configs/config.yaml")
		Conf = &Config{}
		_ = yaml.Unmarshal(configBytes, Conf)
		//设置登录者的缓存
		LoginInfo = &LoginInfoMode{
			Token: "",
			User:  nil,
		}
	})
	return baseConfig.BaseDir, nil
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
