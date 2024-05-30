package api

import (
	api "go-im-service/build/generated/service/v1"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"

	"google.golang.org/protobuf/proto"
)

func IsFriend(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.FriendReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	friendService := service.NewFriendService()
	friend, err := friendService.IsFriend(req.Id)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	if friend != nil {
		return SyncPutSuccess(1, resp)
	}
	return SyncPutSuccess(0, resp)
}

func DelFriend(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.FriendApplyReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewFriendService().DelFriend(req.Id)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func SelectOneFriend(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.FriendReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	f, err := service.NewFriendService().SelectOne(req.Id, false)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(f, resp)
}
func SelectAllFriend() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	fs, err := service.NewFriendService().SelectAll()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(fs, resp)
}
func UpdateFriendName(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.FriendReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewFriendService().UpdateName(req.Id, req.Name)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
