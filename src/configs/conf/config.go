package conf

import (
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"im-sdk/client"
	"os"
	"strings"
	"sync"
)

var (
	DbPath = "im.db"
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
	// 1- 不知道是否需要输入 2-需要输入 -1-不需要输入
	InputPwd2 int
	// 1- 还未输入或者输错了 2-输入正确
	EnteredPwd2 int
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
func InitConfig(baseConfig *BaseConfig) {
	//初始化日志
	once.Do(func() {
		Base = &BaseConfig{
			BaseDir:    baseConfig.BaseDir,
			LogFile:    baseConfig.BaseDir + "/logs",
			ApiHost:    baseConfig.ApiHost,
			WsHost:     baseConfig.WsHost,
			DeviceType: baseConfig.DeviceType,
		}
		log.InitLog(baseConfig.BaseDir+"/logs", baseConfig.LogSwitch)
		//初始化数据库路径
		DbPath = baseConfig.BaseDir + "/configs/db/" + DbPath
		////启动长连接
		c := client.New(Base.WsHost)
		err := c.Startup()
		if err != nil {
			_ = log.WithError(err, "启动长连接失败")
		}
		//读取yaml
		configBytes, _ := os.ReadFile("./config.yaml")
		Conf = &Config{}
		_ = yaml.Unmarshal(configBytes, Conf)
		Conf.ExUris = strings.Split(Conf.Uris, ",")
	})
}
func GetLoginInfo() *LoginInfoMode {
	if LoginInfo == nil {
		LoginInfo = &LoginInfoMode{}
	}
	ReadFile(Base.BaseDir+"/configs/LoginInfo.json", LoginInfo)
	return LoginInfo
}
func PutLoginInfo(user entity.User) {
	LoginInfo = GetLoginInfo()
	LoginInfo.User = &user
	WriteFile(Base.BaseDir+"/configs/LoginInfo.json", LoginInfo)
}
func PutToken(token string) {
	LoginInfo = GetLoginInfo()
	LoginInfo.Token = token
	WriteFile(Base.BaseDir+"/configs/LoginInfo.json", LoginInfo)
}
func UpdateInputPwd2(inputPwd2 int) {
	LoginInfo = GetLoginInfo()
	LoginInfo.InputPwd2 = inputPwd2
	WriteFile(Base.BaseDir+"/configs/LoginInfo.json", LoginInfo)
}
func ReadFile(url string, data any) {
	bytes, _ := os.ReadFile(url)
	json.Unmarshal(bytes, data)
}
func WriteFile(url string, data any) {
	f, _ := os.OpenFile(url, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	//写入
	bytes, _ := json.Marshal(data)
	f.Write(bytes)
}

type Config struct {
	Uris   string
	ExUris []string
	Prime  string
	Pk     string
	Aws    *Aws
}
type Aws struct {
	Id       string
	Secret   string
	Endpoint string
	Region   string
	Bucket   string
}
