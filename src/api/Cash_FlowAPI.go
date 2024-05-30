package api

import (
	api "go-im-service/build/generated/service/v1"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"
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
