package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestSetLanguage(t *testing.T) {
	resp := SetLanguage("en")
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debug(result)
}

func TestSelectConfig(t *testing.T) {
	TestSetLanguage(t)
	resp := SelectConfig()
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debug(result.Msg)
}
func TestGetAgent(t *testing.T) {
	TestSetLanguage(t)
	resp := GetAgent()
	result := &api.ResultDTOResp{}
	proto.Unmarshal(resp, result)
	log.Debug(result.Msg)
}
