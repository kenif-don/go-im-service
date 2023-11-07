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
		Target: 103,
	}
	req, _ := proto.Marshal(oldReq)
	resp := OpenChat(req, nil)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
}
func TestSendMsg(t *testing.T) {
	TestLogin(t)
	oldReq := &api.ChatReq{
		Type:    "friend",
		Target:  103,
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
	time.Sleep(time.Hour * 2)
}
