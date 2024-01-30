package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
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
