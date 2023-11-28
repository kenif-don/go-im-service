package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestGetVersion(t *testing.T) {
	oldReq := &api.VersionReq{
		Type:        1,
		VersionCode: 121,
	}
	req, _ := proto.Marshal(oldReq)
	resp := GetVersion(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
