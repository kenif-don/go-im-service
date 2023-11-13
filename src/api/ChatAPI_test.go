package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestOpenChat(t *testing.T) {
	oldReq := &api.ChatReq{
		Type:   "friend",
		Target: 1,
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
	TestLogin(t)
	TestOpenChat(t)
	contentObj := &entity.MessageData{
		Type:    1,
		Content: "hello",
	}
	content, _ := util.Obj2Str(contentObj)
	oldReq := &api.ChatReq{
		Type:    "friend",
		Target:  1,
		No:      util.Uint642Str(uint64(time.Now().UnixMilli())),
		Content: content,
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
	TestLogin(t)
	resp := GetChats()
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
	time.Sleep(time.Hour)
}
