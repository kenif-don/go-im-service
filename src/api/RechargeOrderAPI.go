package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

// GetRechargeTypes 获取支付类型
func GetRechargeTypes() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	obj := service.NewRechargeOrderService().GetTypes()
	result, e := util.Obj2Str(obj)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = result
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return res
}

// AddRechargeOrder 充值
func AddRechargeOrder(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.RechargeOrderReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	result, err := service.NewRechargeOrderService().AddRechargeOrder(int(req.Type), float64(req.Value))
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = result
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_RECHARGE_FAIL, resp)
	}
	return res
}
