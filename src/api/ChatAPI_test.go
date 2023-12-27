package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestDecrypt(t *testing.T) {
	TestOpenChat(t)
	time.Sleep(time.Second * 2)
	oldReq := &api.DecryptReq{
		Type:   "friend",
		Target: 21,
		No:     "1703499138761",
	}
	req, _ := proto.Marshal(oldReq)
	resp := Decrypt(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
	time.Sleep(time.Hour)
}
func TestOpenChat(t *testing.T) {
	oldReq := &api.ChatReq{
		Type:   "group",
		Target: 22,
	}
	req, _ := proto.Marshal(oldReq)
	resp := OpenChat(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
func TestSendMsg(t *testing.T) {
	TestAutoLogin(t)
	//TestLogin(t)
	time.Sleep(time.Second * 2)
	//TestGetChats(t)
	time.Sleep(time.Second * 2)
	TestOpenChat(t)
	//time.Sleep(time.Second * 2)
	//TestGetMsgs(t)
	contentObj := &api.MessageData{
		//Type:    2,
		//Content: "C:\\Users\\Administrator\\Desktop\\logo.png",
		//Content: "C:\\Users\\Administrator\\Desktop\\b_be930378be7919df8057ce403e1b4d3e.gif",
		Type:    1,
		Content: "成交价啊山莨菪碱扫",
	}
	oldReq := &api.ChatReq{
		Type:    "group",
		Target:  22,
		No:      util.Uint642Str(uint64(time.Now().UnixMilli())),
		Content: contentObj,
	}
	req, _ := proto.Marshal(oldReq)
	resp := SendMsg(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	time.Sleep(time.Hour)
}
func TestGetChats(t *testing.T) {
	resp := GetChats()
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}

func TestGetMsgs(t *testing.T) {
	oldReq := &api.ChatReq{
		Type:   "friend",
		Target: 1,
	}
	req, _ := proto.Marshal(oldReq)
	resp := GetMsgs(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result.Body)
}
func TestDelChatMsg(t *testing.T) {
	TestLogin(t)
	oldReq := &api.ChatReq{
		Type:   "friend",
		Target: 1,
	}
	req, _ := proto.Marshal(oldReq)
	resp := DelChatMsg(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
	time.Sleep(time.Hour)
}
