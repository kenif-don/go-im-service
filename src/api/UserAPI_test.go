package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"fmt"
	"google.golang.org/protobuf/proto"
	"testing"
)

func init() {
	config := &api.ConfigReq{
		BaseDir:    "..",
		LogSwitch:  api.ConfigReq_CONSOLE_FILE,
		DeviceType: api.ConfigReq_Android,
		ApiHost:    "http://hp9kwse9.beesnat.com",
		WsHost:     "ws://ggeejj9f.beesnat.com:13191",
	}
	req, _ := proto.Marshal(config)
	resp := InitConfig(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("配置初始化成功！ %+v", result)
}
func TestRegister(t *testing.T) {
	user := &api.RegisterReq{
		Username: "test123",
		Password: "123456",
	}
	req, _ := proto.Marshal(user)
	resp := Register(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("%+v", result)
}
func TestLogin(t *testing.T) {
	user := &api.RegisterReq{
		Username: "test123",
		Password: "123456",
	}
	req, _ := proto.Marshal(user)
	resp := Login(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("%+v", result)
}

func TestInfo(t *testing.T) {
	resp := Info()
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	fmt.Println(result.Data)
	log.Debugf("%+v", result)
}
