package api

import (
	api "go-im-service/build/generated/service/v1"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"

	"google.golang.org/protobuf/proto"
)

//	type RegisterListener interface {
//		On(data []byte)
//	}
func UpdateFriendApply(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.FriendApplyReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewFriendApplyService().Update(req.Id, int(req.State))
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func SelectFriendApplyNotOperated() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	friendApplyService := service.NewFriendApplyService()
	fas, err := friendApplyService.SelectFriendApplyNotOperated()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(fas, resp)
}
func SelectAllFriendApply() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	friendApplyService := service.NewFriendApplyService()
	fas, err := friendApplyService.SelectAll()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(fas, resp)
}
func AddFriend(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.FriendApplyReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	friendApplyService := service.NewFriendApplyService()
	err := friendApplyService.Add(req.Id, req.Remark)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
