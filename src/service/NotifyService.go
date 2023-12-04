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

func FileNotify(no string, data []byte, state int32, ext string, ext2 []byte) {
	if Listener != nil {
		fr := &api.FileDecryptResp{
			No:    no,
			Data:  data,
			State: state,
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
