package api

import (
	api "go-im-service/build/generated/service/v1"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"

	"google.golang.org/protobuf/proto"
)

// GetRechargeTypes 获取支付类型
func GetRechargeTypes() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	obj := service.NewRechargeOrderService().GetTypes()
	return SyncPutSuccess(obj, resp)
}

// AddRechargeOrder 充值
func AddRechargeOrder(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.RechargeOrderReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	result, err := service.NewRechargeOrderService().AddRechargeOrder(int(req.Type), req.Value)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(result, resp)
}
