package service

import (
	"IM-Service/src/configs/db"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/util"
	"im-sdk/model"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

// GetOfflineMessage 获取离线消息
func (_self *MessageService) GetOfflineMessage() *utils.Error {
	resultDTO, err := util.Post("/api/offline-bill/selectAll", nil)
	if err != nil {
		return log.WithError(err)
	}
	var offlineBills = &[]entity.OfflineBill{}
	e := util.Str2Obj(resultDTO.Data.(string), offlineBills)
	if e != nil {
		return log.WithError(utils.ERR_QUERY_FAIL)
	}
	tx := db.NewTransaction().BeginTx()
	if err := tx.Error; err != nil {
		return log.WithError(utils.ERR_QUERY_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = func() *utils.Error {
		var ids = make([]int64, len(*offlineBills))
		for i, bill := range *offlineBills {
			var protocol = &model.Protocol{}
			e = util.Str2Obj(bill.Content, protocol)
			if e != nil {
				return log.WithError(utils.ERR_QUERY_FAIL)
			}
			ids[i] = bill.Id
			err := _self.Handler(protocol)
			if err != nil {
				return log.WithError(err)
			}
		}
		if len(ids) == 0 {
			return nil
		}
		var req = make(map[string]interface{})
		req["ids"] = ids
		_, err = util.Post("/api/offline-bill/dels", req)
		if err != nil {
			return log.WithError(err)
		}
		return nil
	}()
	return err
}
func (_self *MessageService) Handler(protocol *model.Protocol) *utils.Error {
	switch protocol.Type {
	case 101: //让to拉去from的好友申请信息，没有就存起来 有就修改
		err := NewFriendApplyService().updateOne(util.Str2Uint64(protocol.From), util.Str2Uint64(protocol.To))
		if err != nil {
			return log.WithError(err)
		}
		break
	case 102: //当to同意好友申请后，更新好友数据
		err := NewFriendService().updateOne(util.Str2Uint64(protocol.From), util.Str2Uint64(protocol.To))
		if err != nil {
			return log.WithError(err)
		}
		break
	}
	return nil
}
