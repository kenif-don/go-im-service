package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestAddFriend(t *testing.T) {
	oldReq := &api.FriendApplyReq{
		Id:     2,
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
	resp := SelectAllFriendApply()
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(result)
}
func TestUpdateFriendApply(t *testing.T) {
	oldReq := &api.FriendApplyReq{
		Id:    1,
		State: -1,
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
