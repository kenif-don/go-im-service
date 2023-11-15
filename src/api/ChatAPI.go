package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

// DelChatMsg 删除双方聊天消息
func DelChatMsg(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	messageService := service.NewMessageService()
	err := messageService.DelChatMsg(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return res
}

// DelLocalChatMsg 删除聊天消息
func DelLocalChatMsg(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	messageService := service.NewMessageService()
	err := messageService.DelLocalChatMsg(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return res
}
func GetConnectState() []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	resp := &api.ResultDTOResp{}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	if conf.Conf.Connected {
		resp.Body = "1"
	} else {
		resp.Body = "0"
	}
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	return res
}

// DelLocalChat 删除本地聊天记录
func DelLocalChat(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	chatService := service.NewChatService()
	err := chatService.DelLocalChat(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_DEL_FAIL, resp)
	}
	return res
}

// DelChat 删除双方聊天记录
func DelChat(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	chatService := service.NewChatService()
	err := chatService.DelChat(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_DEL_FAIL, resp)
	}
	return res
}
func GetChats() []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	resp := &api.ResultDTOResp{}
	chatService := service.NewChatService()
	chats, err := chatService.GetChats()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	result, e := util.Obj2Str(chats)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = result
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_OPEN_FAIL, resp)
	}
	return res
}
func OpenChat(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	chatService := service.NewChatService()
	chats, err := chatService.OpenChat(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	result, e := util.Obj2Str(chats)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = result
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_OPEN_FAIL, resp)
	}
	return res
}
func SendMsg(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	msgService := service.NewMessageService()
	err := msgService.SendMsg(req.Type, req.Target, req.No, req.Content)
	if err != nil {
		return SyncPutErr(utils.ERR_SEND_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_SEND_FAIL, resp)
	}
	return res
}

// GetMsgs 分页获取消息
func GetMsgs(data []byte) []byte {
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, nil)
	}
	req := &api.MsgPageReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	messageService := service.NewMessageService()
	msgs, err := messageService.Paging(req.Type, req.Target, req.Time)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	result, e := util.Obj2Str(msgs)
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
