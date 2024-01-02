package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

func GetGroup(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	group, err := service.NewGroupService().SelectOne(req.Id, false)
	if err != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	obj, e := util.Obj2Str(group)
	if e != nil {
		return SyncPutErr(utils.ERR_OPERATION_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = obj
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_ADD_MOOD_FAIL, resp)
	}
	return res
}
func GetGroupMembers(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	members, err := service.NewGroupMemberService().SelectMembers(req.Id, false)
	if err != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	obj, e := util.Obj2Str(members)
	if e != nil {
		return SyncPutErr(utils.ERR_OPERATION_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = obj
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_ADD_MOOD_FAIL, resp)
	}
	return res
}

// GetGroups 获取群
func GetGroups() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	groups, e := service.NewGroupService().SelectAll()
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	obj, e := util.Obj2Str(groups)
	if e != nil {
		return SyncPutErr(utils.ERR_OPERATION_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = obj
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_ADD_MOOD_FAIL, resp)
	}
	return res
}

// CreateGroup 创建群聊
func CreateGroup(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	group, err := service.NewGroupService().Create(util.Str2Arr(req.Ids), int(req.Type), req.Password)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	obj, e := util.Obj2Str(group)
	if e != nil {
		return SyncPutErr(utils.ERR_OPERATION_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = obj
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_ADD_MOOD_FAIL, resp)
	}
	return res
}

// InviteInGroup 邀请进群
func InviteInGroup(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewGroupService().Invite(req.Id, util.Str2Arr(req.Ids))
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_ADD_MOOD_FAIL, resp)
	}
	return res
}
