package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestGetWithdrawalFee(t *testing.T) {
	resp := GetWithdrawalFee()
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
