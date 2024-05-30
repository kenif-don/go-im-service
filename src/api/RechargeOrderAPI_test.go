package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"

	"google.golang.org/protobuf/proto"
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
