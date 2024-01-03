package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"github.com/go-netty/go-netty-transport/websocket"
	"google.golang.org/protobuf/proto"
)

// ValidateGroupNeedPassword 验证群聊是否需要密码
func ValidateGroupNeedPassword(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.ChatReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	obj := service.NewGroupService().NeedPassword(req.Type, req.Target)
	return SyncPutSuccess(obj, resp)
}
func Decrypt(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.DecryptReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.DecryptFile(req.Type, req.Target, req.No)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// DelChatMsg 删除双方聊天消息
func DelChatMsg(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.ChatReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	messageService := service.NewMessageService()
	err := messageService.DelChatMsg(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// DelLocalChatMsg 删除聊天消息
func DelLocalChatMsg(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.ChatReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	messageService := service.NewMessageService()
	err := messageService.DelLocalChatMsg(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func ImReConnect() []byte {
	resp := &api.ResultDTOResp{}
	e := conf.Conf.Client.Reconnect(websocket.New())
	if e == nil {
		return SyncPutErr(utils.ERR_NET_FAIL, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// DelLocalChat 删除本地聊天记录
func DelLocalChat(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.ChatReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	chatService := service.NewChatService()
	err := chatService.DelLocalChat(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// DelChat 删除双方聊天记录
func DelChat(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.ChatReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	chatService := service.NewChatService()
	err := chatService.DelChat(req.Type, req.Target)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
func GetChats() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	chatService := service.NewChatService()
	chats, err := chatService.GetChats()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(chats, resp)
}
func OpenChat(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.ChatReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	chat, err := service.NewChatService().OpenChat(req.Type, req.Target, req.Password)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(chat, resp)
}
func SendMsg(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.ChatReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	msgService := service.NewMessageService()
	go func() {
		err := msgService.SendMsg(req.Type, req.Target, req.No, &entity.MessageData{
			Type:    int(req.Content.Type),
			Content: req.Content.Content,
		})
		if err != nil && service.Listener != nil {
			//通知消息发送失败
			err = service.NotifySendReceive(req.No, -1)
			if err != nil {
				log.Error(err)
			}
		}
	}()
	return SyncPutSuccess(nil, resp)
}

// GetMsgs 分页获取消息
func GetMsgs(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.MsgPageReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	messageService := service.NewMessageService()
	msgs, err := messageService.Paging(req.Type, req.Target, req.Time)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(msgs, resp)
}
func CurrentTime() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	t := util.Uint642Str(service.NewMessageService().CurrentTime())
	return SyncPutSuccess(t, resp)
}
