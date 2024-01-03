package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
)

// SelectRemoteAccount 获取登录者账户数据 没有就从服务器同步
func SelectRemoteAccount() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	obj, err := service.NewAccountService().SelectOneAccount(true)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(obj, resp)
}

// SelectOneAccount 从本地获取账户数据
func SelectOneAccount() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	obj, err := service.NewAccountService().SelectOneAccount(false)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(obj, resp)
}
