package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

func DelLocalChat(data []byte) []byte {
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
func DelChat(data []byte) []byte {
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
