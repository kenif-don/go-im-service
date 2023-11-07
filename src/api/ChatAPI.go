package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

type MessageListener interface {
	//OnReceive 当前聊天接收到消息
	OnReceive(data []byte)
	//OnSendReceive 发送的消息状态 -某消息 发送成功、发送失败
	OnSendReceive(data []byte)
}

func OpenChat(data []byte, listener MessageListener) []byte {
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	chatService := service.NewChatService(listener)
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
	resp.Data = result
	res, _ := proto.Marshal(resp)
	return res
}
func SendMsg(data []byte) []byte {
	req := &api.ChatReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	msgService := service.NewMessageService()
	msgService.SendMsg(req.Type, req.Target, req.No, req.Content)
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, _ := proto.Marshal(resp)
	return res
}
