package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestAddMood(t *testing.T) {
	oldReq := &api.MoodReq{
		Content: "test",
		Urls:    "['https://www.baidu.com']",
	}
	req, _ := proto.Marshal(oldReq)
	resp := AddMood(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}

func TestDeleteMood(t *testing.T) {
	oldReq := &api.MoodReq{
		Id: 36,
	}
	req, _ := proto.Marshal(oldReq)
	resp := DeleteMood(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}

func TestSelectOneMood(t *testing.T) {
	oldReq := &api.MoodReq{
		Id: 37,
	}
	req, _ := proto.Marshal(oldReq)
	resp := SelectOneMood(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}

func TestPagingMood(t *testing.T) {
	oldReq := &api.MoodPageReq{
		Page:     1,
		PageSize: 10,
		//UserId:   1,
	}
	req, _ := proto.Marshal(oldReq)
	resp := PagingMood(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}

func TestAddReply(t *testing.T) {
	oldReq := &api.ReplyReq{
		MoodId: 37,
		//ReplyUserId: 5,
		Content: "test",
	}
	req, _ := proto.Marshal(oldReq)
	resp := AddReply(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
