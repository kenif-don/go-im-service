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
		Value: "7.7",
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

func TestAddWithdrawal(t *testing.T) {
	oldReq := &api.WithdrawalReq{
		Money:   "100",
		Address: "TGm6v1BFdCnfWygtvYU4wp2EXMRuWbWuYo",
	}
	req, _ := proto.Marshal(oldReq)
	resp := AddWithdrawal(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
