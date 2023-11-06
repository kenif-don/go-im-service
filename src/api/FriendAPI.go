package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

func DelFriend(data []byte) []byte {
	req := &api.FriendApplyReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	friendService := service.NewFriendService()
	err := friendService.Del(req.Id)
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
func SelectOneFriend(data []byte) []byte {
	req := &api.FriendReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	friendService := service.NewFriendService()
	f, err := friendService.SelectOne(req.Id)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	result, e := util.Obj2Str(f)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Data = result
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return res
}
func SelectAllFriend() []byte {
	resp := &api.ResultDTOResp{}
	friendService := service.NewFriendService()
	fs, err := friendService.SelectAll()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	result, e := util.Obj2Str(fs)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Data = result
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return res
}
func UpdateFriendName(data []byte) []byte {
	req := &api.FriendReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	friendService := service.NewFriendService()
	err := friendService.UpdateName(req.Id, req.Name)
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
