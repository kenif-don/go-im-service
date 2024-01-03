package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestGetGroupMemberInfo(t *testing.T) {
	TestAutoLogin(t)
	time.Sleep(time.Second)
	oldReq := &api.GroupMemberInfoReq{
		GId:    44,
		UserId: 10,
	}
	req, _ := proto.Marshal(oldReq)
	resp := GetGroupMemberInfo(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result.Body)
	time.Sleep(time.Hour)
}
