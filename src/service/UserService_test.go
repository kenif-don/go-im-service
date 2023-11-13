package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/log"
	"testing"
)

func init() {
	conf.InitConfig(&conf.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		DeviceType: "PC",
		ApiHost:    "http://hp9kwse9.beesnat.com",
		WsHost:     "ws://ggeejj9f.beesnat.com:13191",
	})
	log.Debugf("配置初始化成功!")
}
func TestUserService_Register(t *testing.T) {
	userService := NewUserService()
	err := userService.Register("test", "123456")
	if err != nil {
		log.Debug(err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	NewUserService().UpdateUser(7)
}
