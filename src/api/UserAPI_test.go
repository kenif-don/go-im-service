package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/log"
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
	resp := InitConfig(req, nil)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("配置初始化成功！ %+v", result)
}
func TestRegister(t *testing.T) {
	user := &api.UserReq{
		//Username: "test123",
		//Password: "123456",
		//Username: "666666",
		//Password: "666666",
		Username: "66666666666666",
		Password: "66666666666666",
	}
	req, _ := proto.Marshal(user)
	resp := Register(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("%+v", result)
}
func TestLogin(t *testing.T) {
	user := &api.UserReq{
		//Username: "1338076457770",
		//Password: "a12345678",
		Username: "test123",
		Password: "123456",
		//Username: "123456",
		//Password: "123456",
	}
	req, _ := proto.Marshal(user)
	resp := Login(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debug(result)
}

func TestInfo(t *testing.T) {
	resp := Info()
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debug(result)
}

func TestUpdateNickname(t *testing.T) {
	updateReq := &api.UpdateUserReq{
		Id:   conf.GetLoginInfo().User.Id,
		Data: "冷风",
	}
	req, _ := proto.Marshal(updateReq)
	resp := UpdateNickname(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("%+v", result)
}
func TestUpdateIntro(t *testing.T) {
	updateReq := &api.UpdateUserReq{
		Id:   conf.GetLoginInfo().User.Id,
		Data: "冷霜自拌无情雨，孤叶何须罪秋风",
	}
	req, _ := proto.Marshal(updateReq)
	resp := UpdateIntro(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("%+v", result)
}
func TestUpdateEmail(t *testing.T) {
	updateReq := &api.UpdateUserReq{
		Id:   conf.GetLoginInfo().User.Id,
		Data: "123@qq.com",
	}
	req, _ := proto.Marshal(updateReq)
	resp := UpdateEmail(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("%+v", result)
}
func TestUpdateHeadImg(t *testing.T) {
	updateReq := &api.UpdateUserReq{
		Id:   conf.GetLoginInfo().User.Id,
		Data: "https://123.png",
	}
	req, _ := proto.Marshal(updateReq)
	resp := UpdateHeadImg(req)
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debugf("%+v", result)
}
func TestSearch(t *testing.T) {
	searchReq := &api.SearchReq{
		Keyword: "冷风",
	}
	req, _ := proto.Marshal(searchReq)
	resp := Search(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(result)
}
func TestLogout(t *testing.T) {
	resp := Logout()
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debug(result)
}
