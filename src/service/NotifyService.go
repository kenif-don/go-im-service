package service

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
	"time"
)

// OfflineMessageNotify 获取离线消息
func OfflineMessageNotify() {
	go func() {
		for {
			//如果没有私钥 或者还未登录 就下一次循环  TODO 这里可能会在登录之前调用 而且秘钥可能有问题 所以离线消息会解密失败 成为乱码
			if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.PrivateKey == "" {
				time.Sleep(time.Second * 1)
				continue
			}
			//获取离线消息
			err := NewMessageService().GetOfflineMessage()
			if err != nil {
				log.Error(err)
			}
			return
		}
	}()
}
func FileNotify(target uint64, no, content string) *utils.Error {
	if conf.Conf.ChatId != target {
		return nil
	}
	if Listener != nil {
		//根据No获取消息记录
		message, err := NewMessageService().SelectOne(&entity.Message{No: no})
		if err != nil {
			return log.WithError(err)
		}
		resp := &api.FileDecryptResp{
			No:      no,
			Content: content,
			Type:    message.Type,
		}
		res, e := proto.Marshal(resp)
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_NOTIFY_FAIL)
		}
		Listener.OnFile(res)
	}
	return nil
}
func DelMsgNotify(tp string, target uint64) *utils.Error {
	if conf.Conf.ChatId != target {
		return nil
	}
	if Listener != nil {
		chat := &entity.Chat{
			Type:     tp,
			TargetId: target,
		}
		data, e := util.Obj2Str(chat)
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_NOTIFY_FAIL)
		}
		Listener.OnDelMsg(data)
	}
	return nil
}

// ChatNotify 通知客户端更新聊天列表
func (_self *ChatService) ChatNotify(chat *entity.Chat) *utils.Error {
	err := _self.coverLastMsg(chat)
	if err != nil {
		return log.WithError(err)
	}
	if Listener != nil {
		res, e := util.Obj2Str(chat)
		if e != nil {
			return log.WithError(e)
		}
		Listener.OnDoChat(res)
	}
	return nil
}

func (_self *ChatService) VoiceNotify(message *entity.Message) *utils.Error {
	if Listener != nil {
		res, e := util.Obj2Str(message)
		if e != nil {
			return log.WithError(e)
		}
		Listener.OnDoVoice(res)
	}
	return nil
}

// NotifySendReceive 通知消息是否发生成功
func NotifySendReceive(no string, send int) *utils.Error {
	if no == "" {
		return nil
	}
	//根据No获取消息记录
	message, err := NewMessageService().SelectOne(&entity.Message{No: no})
	if err != nil {
		return log.WithError(err)
	}
	//修改消息状态
	if message != nil {
		message.Send = send
		err := NewMessageService().Update(message)
		if err != nil {
			return log.WithError(err)
		}
	}
	if Listener != nil {
		m := map[string]interface{}{"no": no, "send": send, "type": 0}
		if message != nil {
			m["type"] = message.Type
		}
		res, e := util.Obj2Str(m)
		if e != nil {
			log.Error(e)
			return log.WithError(e)
		}
		Listener.OnSendReceive(res)
	}
	return nil
}
func NotifyReceive(message *entity.Message, target uint64) *utils.Error {
	if message.TargetId == conf.Conf.ChatId {
		//解密
		data, err := Decrypt(message.Type, target, message.No, message.Data)
		if err != nil {
			data = util.GetTextErrMsg()
		}
		message.Data = data
		if Listener != nil {
			res, e := util.Obj2Str(message)
			if e != nil {
				log.Error(e)
				return log.WithError(e)
			}
			//修改聊天消息为已读
			NewMessageService().UpdateChatRead(message.Type, target)
			Listener.OnReceive(res)
		}
	}
	return nil
}
