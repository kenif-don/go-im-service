package api

import (
	api "go-im-service/build/generated/service/v1"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"

	"google.golang.org/protobuf/proto"
)

// CreateRedeemCode 钱兑码
func CreateRedeemCode(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.RedeemCodeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	code, err := service.NewRedeemCodeService().Create(req.Money)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(code, resp)
}

// Exchange 码兑钱
func Exchange(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.RedeemCodeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewRedeemCodeService().Exchange(req.Code)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// PagingRedeemCode 分页获取兑换记录
func PagingRedeemCode(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.RedeemCodeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	code, err := service.NewRedeemCodeService().Paging(int(req.Page), int(req.PageSize))
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(code, resp)
}

// SelectOneRedeemCode 获取兑换记录
func SelectOneRedeemCode(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.RedeemCodeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	code, err := service.NewRedeemCodeService().SelectOne(req.Code)
	if err != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return SyncPutSuccess(code, resp)
}
