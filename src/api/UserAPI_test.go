package api

//
//import (
//	"IM-Service/configs/log"
//	api "IM-Service/generated/service/v1"
//	"google.golang.org/protobuf/proto"
//	"testing"
//)
//
//func init() {
//	config := &api.ConfigReq{
//		BaseDir:    "..",
//		LogSwitch:  api.ConfigReq_CONSOLE_FILE,
//		DeviceType: api.ConfigReq_Android,
//		ApiHost:    "http://127.0.0.1:8886",
//		WsHost:     "ws://127.0.0.1:8003",
//	}
//	req, _ := proto.Marshal(config)
//	resp := InitConfig(req)
//	result := &api.ResultDTOResp{}
//	proto.Unmarshal(resp, result)
//	log.Debugf("%+v", result)
//}
//func TestRegister(t *testing.T) {
//	user := &api.RegisterReq{
//		Username: "test123",
//		Password: "123456",
//	}
//	req, _ := proto.Marshal(user)
//	resp := Register(req)
//	result := &api.ResultDTOResp{}
//	proto.Unmarshal(resp, result)
//	log.Debugf("%+v", result)
//}
