package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestAddRechargeOrder(t *testing.T) {
	oldReq := &api.RechargeOrderReq{
		Type:  1,
		Value: "1",
	}
	req, _ := proto.Marshal(oldReq)
	resp := AddRechargeOrder(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
