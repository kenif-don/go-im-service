package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestIsFriend(t *testing.T) {
	resp := IsFriend(nil)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
func TestDelFriend(t *testing.T) {
	TestLogin(t)
	oldReq := &api.FriendReq{
		Id: 144,
	}
	req, _ := proto.Marshal(oldReq)
	resp := DelFriend(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
func TestSelectAllFriend(t *testing.T) {
	TestLogin(t)
	resp := SelectAllFriend()
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
func TestSelectOneFriend(t *testing.T) {
	oldReq := &api.FriendReq{
		Id: 103,
	}
	req, _ := proto.Marshal(oldReq)
	resp := SelectOneFriend(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
func TestUpdateName(t *testing.T) {
	TestLogin(t)
	oldReq := &api.FriendReq{
		Id:   103,
		Name: "测试人员",
	}
	req, _ := proto.Marshal(oldReq)
	resp := UpdateFriendName(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
}
