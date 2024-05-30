package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestAddFriend(t *testing.T) {
	TestAutoLogin(t)
	oldReq := &api.FriendApplyReq{
		Id:     21,
		Remark: "加我",
	}
	req, _ := proto.Marshal(oldReq)
	resp := AddFriend(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(result)
}
func TestSelectAllFriendApply(t *testing.T) {
	TestLogin(t)
	resp := SelectAllFriendApply()
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Debug(err)
	}

	log.Debug(result)
	time.Sleep(time.Hour)
}
func TestUpdateFriendApply(t *testing.T) {
	oldReq := &api.FriendApplyReq{
		Id:    8,
		State: 2,
	}
	req, _ := proto.Marshal(oldReq)
	resp := UpdateFriendApply(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(result)
}
