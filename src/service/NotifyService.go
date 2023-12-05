package service

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

func FileNotify(no, path string, state, tp int32, ext string, ext2 []byte) {
	if Listener != nil {
		fr := &api.FileDecryptResp{
			No:    no,
			Path:  path,
			State: state,
			Type:  tp,
			Ext:   ext,
			Ext2:  ext2,
		}
		res, e := proto.Marshal(fr)
		if e != nil {
			log.Debug(e)
		}
		Listener.OnFile(res)
	}
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
