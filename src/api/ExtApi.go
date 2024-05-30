package api

import (
	api "go-im-service/build/generated/service/v1"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"
)

func GetExtUrl() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	ext, err := service.NewExtService().Get()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(ext, resp)
}
