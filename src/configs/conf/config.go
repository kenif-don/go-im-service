package conf

import (
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"encoding/json"
	"im-sdk/client"
	"os"
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
	DeviceType string
}

// InitConfig 初始化方法
func InitConfig(baseConfig *BaseConfig) {
	//初始化日志
	once.Do(func() {
		Base = &BaseConfig{
			BaseDir:    baseConfig.BaseDir,
			ApiHost:    baseConfig.ApiHost,
			WsHost:     baseConfig.WsHost,
			DeviceType: baseConfig.DeviceType,
		}
		log.InitLog(mkdirAndReturn(baseConfig.BaseDir+"/logs"), baseConfig.LogSwitch)
		//初始化数据库路径
		DbPath = mkdirAndReturn(baseConfig.BaseDir+"/configs/db") + "/" + DbPath
		log.Debug("数据库路径:", DbPath)
		//读取yaml
		Conf = &Config{
			ExUris: []string{"/api/user/login", "/api/user/info", "/api/user/register", "/api/user/resetPublicKey", "/api/test/index", "/back/admin/login", "/back/admin/info", "/back/admin/resetPublicKe"},
			Prime:  "262074f1e0e19618f0d2af786779d6ad9e814b",
			Pk:     "19311a1a18656914b9381c058c309083022301",
			Aws: &Aws{
				Id:       "WHZGIHUSERKPOCCITVOPDJPM",
				Secret:   "haYlDZAdsSN4zckmX64W0zKSDg7IWxdH1lOkxm9N",
				Endpoint: "https://world.sgs3.layerstackobjects.com",
				Region:   "Singapore",
				Bucket:   "world",
			},
		}
	})
}
func ClearLoginInfo() {
	LoginInfo = nil
	WriteFile(mkdirAndReturn(Base.BaseDir+"/configs")+"/LoginInfo.json", &LoginInfoMode{})
}
func GetLoginInfo() *LoginInfoMode {
	if LoginInfo == nil {
		LoginInfo = &LoginInfoMode{}
	}
	ReadFile(mkdirAndReturn(Base.BaseDir+"/configs")+"/LoginInfo.json", LoginInfo)
	return LoginInfo
}
func PutLoginInfo(user entity.User) {
	LoginInfo = GetLoginInfo()
	LoginInfo.User = &user
	WriteFile(mkdirAndReturn(Base.BaseDir+"/configs")+"/LoginInfo.json", LoginInfo)
}
func PutToken(token string) {
	LoginInfo = GetLoginInfo()
	LoginInfo.Token = token
	WriteFile(mkdirAndReturn(Base.BaseDir+"/configs")+"/LoginInfo.json", LoginInfo)
}
func UpdateInputPwd2(inputPwd2 int) {
	LoginInfo = GetLoginInfo()
	LoginInfo.InputPwd2 = inputPwd2
	WriteFile(mkdirAndReturn(Base.BaseDir+"/configs")+"/LoginInfo.json", LoginInfo)
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
func mkdirAndReturn(path string) string {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	return path
}

type Config struct {
	ExUris    []string
	Prime     string
	Pk        string
	Aws       *Aws
	Key       string // 与服务器交互的key
	Client    *client.WSClient
	Connected bool   //长连接是否链接成功
	ChatId    uint64 //当前打开的聊天ID
}
type Aws struct {
	Id       string
	Secret   string
	Endpoint string
	Region   string
	Bucket   string
}
