package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestAddFriend(t *testing.T) {

	oldReq := &api.FriendApplyReq{
		Id:     1,
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
