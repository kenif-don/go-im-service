package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
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
	select {
	case <-time.After(time.Second * 35):
		log.Debug("测试超时")
	}
}
