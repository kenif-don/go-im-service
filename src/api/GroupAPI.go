package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

func GetGroupMemberByMe(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	members, err := service.NewGroupMemberService().GetGroupMemberByMe(req.Id)
	if err != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return SyncPutSuccess(members, resp)
}
func GetGroupMemberInfo(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupMemberInfoReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	f, err := service.NewGroupService().SelectOneGroupMemberInfo(req.GId, req.UserId)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(f, resp)
}
func DeleteGroup(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewGroupService().Delete(req.Id)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func QuitGroup(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewGroupService().Quit(req.Id)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func UpdateGroupHeadImg(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewGroupService().Update(req.Id, req.Data, 3)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func UpdateGroupNotice(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewGroupService().Update(req.Id, req.Data, 2)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func UpdateGroupName(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.GroupReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewGroupService().Update(req.Id, req.Data, 1)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func SelectOneGroup(data []byte) []byte {
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
	return SyncPutSuccess(group, resp)
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
	return SyncPutSuccess(members, resp)
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
	return SyncPutSuccess(groups, resp)
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
	return SyncPutSuccess(group, resp)
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
	return SyncPutSuccess(nil, resp)
}
