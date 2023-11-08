package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestOpenChat(t *testing.T) {
	oldReq := &api.ChatReq{
		Type:   "friend",
		Target: 2,
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
	oldReq := &api.ChatReq{
		Type:    "friend",
		Target:  2,
		No:      "2",
		Content: "hello",
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
	TestOpenChat(t)
	resp := GetChats()
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
