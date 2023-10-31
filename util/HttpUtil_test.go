package util

import (
	"IM-Service/configs/conf"
	"IM-Service/configs/log"
	"IM-Service/entity"
	"fmt"
	"testing"
)

func init() {
	conf.InitConfig(&conf.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		DeviceType: "UNKNOWN",
		ApiHost:    "http://127.0.0.1:8886",
		WsHost:     "ws://127.0.0.1:8003",
	})
}
func TestPost(t *testing.T) {
	params := &entity.RegisterUser{
		Username: "test123",
		Password: "123456",
	}
	result, err := Post("/api/user/register", params)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(result.Code, result.Msg, result.Data)
}
