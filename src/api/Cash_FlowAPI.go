package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
)

// PagingCashFlow 分页查询资金流水
func PagingCashFlow(page int) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	obj, err := service.NewCashFlowService().Paging(page)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(obj, resp)
}
