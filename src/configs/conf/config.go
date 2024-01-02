package conf

import (
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"encoding/json"
	"im-sdk/client"
	"os"
	"path/filepath"
	"sync"
)

const (
	PC      = "PC"
	Android = "Android"
	IOS     = "IOS"
)

var (
	DbPath = "im.db"
)
var (
	once      sync.Once
	Base      *BaseConfig
	Conf      *Config
	LoginInfo *LoginInfoMode
	DiffTime  int
)

type LoginInfoMode struct {
	Token string
	User  *entity.User
	// 1- 需要输入 2-已正确输入 -1-不需要输入
	InputPwd2 int
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
		//获取文件信息
		path, err := filepath.Abs(baseConfig.BaseDir)
		if err != nil {
			return
		}
		log.InitLog(mkdirAndReturn(filepath.Join(path, "logs")), baseConfig.LogSwitch)
		//删除临时文件
		e := os.RemoveAll(mkdirAndReturn(filepath.Join(path, "configs", "temp")))
		if e != nil {
			log.Error(e)
		}
		//创建一个空文件夹
		mkdirAndReturn(filepath.Join(path, "configs", "temp"))
		Base = &BaseConfig{
			BaseDir:    path,
			ApiHost:    baseConfig.ApiHost,
			WsHost:     baseConfig.WsHost,
			DeviceType: baseConfig.DeviceType,
		}
		//初始化数据库路径
		DbPath = filepath.Join(mkdirAndReturn(filepath.Join(path, "configs", "db")), DbPath)
		//初始化配置
		Conf = &Config{
			ExUris: []string{"/api/user/login", "/api/user/info", "/api/user/register", "/api/user/resetPublicKey",
				"/api/test/index", "/back/admin/login", "/back/admin/info", "/back/admin/resetPublicKe", "/api/version/select"},
			Prime: "262074f1e0e19618f0d2af786779d6ad9e814b",
			Pk:    "19311a1a18656914b9381c058c309083022301",
			Aws: &Aws{
				Id:       "WHZGIHUSERKPOCCITVOPDJPM",
				Secret:   "haYlDZAdsSN4zckmX64W0zKSDg7IWxdH1lOkxm9N",
				Endpoint: "sgs3.layerstackobjects.com",
				Region:   "Singapore",
				Bucket:   "world",
			},
		}
		//初始化时 判断是否需要二级密码
		if GetLoginInfo().User == nil {
			return
		}
		if GetLoginInfo().User.Password2 != "" {
			//需要输入二级密码
			UpdateInputPwd2(1)
		} else {
			//不需要输入二级密码
			UpdateInputPwd2(-1)
		}
		log.Debug("配置初始化完成")
	})
}
func ClearLoginInfo() {
	LoginInfo = nil
	WriteFile(filepath.Join(mkdirAndReturn(filepath.Join(Base.BaseDir, "configs")), "LoginInfo.json"), &LoginInfoMode{})
}
func GetLoginInfo() *LoginInfoMode {
	if LoginInfo == nil {
		LoginInfo = &LoginInfoMode{}
	}
	ReadFile(filepath.Join(mkdirAndReturn(filepath.Join(Base.BaseDir, "configs")), "LoginInfo.json"), LoginInfo)
	return LoginInfo
}
func PutLoginInfo(user entity.User) {
	LoginInfo = GetLoginInfo()
	LoginInfo.User = &user
	WriteFile(filepath.Join(mkdirAndReturn(filepath.Join(Base.BaseDir, "configs")), "LoginInfo.json"), LoginInfo)
}
func PutToken(token string) {
	LoginInfo = GetLoginInfo()
	LoginInfo.Token = token
	WriteFile(filepath.Join(mkdirAndReturn(filepath.Join(Base.BaseDir, "configs")), "LoginInfo.json"), LoginInfo)
}
func UpdateInputPwd2(inputPwd2 int) {
	LoginInfo = GetLoginInfo()
	LoginInfo.InputPwd2 = inputPwd2
	WriteFile(filepath.Join(mkdirAndReturn(filepath.Join(Base.BaseDir, "configs")), "LoginInfo.json"), LoginInfo)
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
	Pk        string //服务器公钥
	Aws       *Aws
	Key       string // 与服务器交互的key
	Client    *client.WSClient
	Connected bool              //长连接是否链接成功
	ChatId    uint64            //当前打开的聊天ID
	Pwds      map[string]string //密聊群的密码
}
type Aws struct {
	Id       string
	Secret   string
	Endpoint string
	Region   string
	Bucket   string
}
