package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
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

func TestGetGroupMembers(t *testing.T) {
	TestAutoLogin(t)
	//time.Sleep(time.Second * 5)
	//oldReq := &api.GroupReq{
	//	Id: 24,
	//}
	//req, _ := proto.Marshal(oldReq)
	//resp := GetGroupMembers(req)
	//result := &api.ResultDTOResp{}
	//err := proto.Unmarshal(resp, result)
	//if err != nil {
	//	log.Error(err)
	//}
	//log.Debug(result.Body)
	time.Sleep(time.Hour)
}
