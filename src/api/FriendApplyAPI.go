package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

//	type RegisterListener interface {
//		On(data []byte)
//	}
func UpdateFriendApply(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.FriendApplyReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	friendApplyService := service.NewFriendApplyService()
	err := friendApplyService.Update(req.Id, int(req.State))
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
func SelectFriendApplyNotOperated() []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	resp := &api.ResultDTOResp{}
	friendApplyService := service.NewFriendApplyService()
	fas, err := friendApplyService.SelectFriendApplyNotOperated()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	result, e := util.Obj2Str(fas)
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
func SelectAllFriendApply() []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	resp := &api.ResultDTOResp{}
	friendApplyService := service.NewFriendApplyService()
	fas, err := friendApplyService.SelectAll()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	result, e := util.Obj2Str(fas)
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
func AddFriend(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.FriendApplyReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	friendApplyService := service.NewFriendApplyService()
	err := friendApplyService.Add(req.Id, req.Remark)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
