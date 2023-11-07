package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/db"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"im-sdk/handler"
	"im-sdk/model"
	"strconv"
	"time"
)

type MessageService struct {
	repo *repository.MessageRepo
}

func NewMessageService() *MessageService {
	return &MessageService{
		repo: repository.NewMessageRepo(),
	}
}
func (_self *MessageService) QueryLast(obj *entity.Message) (*entity.Message, error) {
	return _self.repo.QueryLast(obj)
}
func (_self *MessageService) UpdateReaded(protocol *model.Protocol, ext4 int) {
	var message = &entity.Message{}
	e := util.Str2Obj(protocol.Data.(string), message)
	if e != nil {
		log.Error(e)
		return
	}
	//修改消息发送状态
	message.Ext4 = ext4
	e = _self.repo.Save(message)
	if e != nil {
		log.Error(e)
	}
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
	case 1: // 接收到聊天消息
		switch protocol.Type {
		case model.ChannelOne2oneMsg, model.ChannelGroupMsg:
			log.Debug("相对===================================8")
			break
		}
		break
	}
	return nil
}

func (_self *MessageService) SendMsg(tp string, target uint64, no, content string) *utils.Error {
	//组装本地消息
	message := _self.coverMessage(tp, target, no, content)
	//组装长连接protocol
	protocol, err := _self.coverProtocol(message)
	if err != nil {
		return log.WithError(err)
	}
	//发送消息
	handler.GetClientHandler().GetMessageManager().Send(protocol)
	//报错消息到数据库
	e := _self.repo.Save(message)
	if e != nil {
		return log.WithError(utils.ERR_SEND_FAIL)
	}
	return nil
}
func (_self *MessageService) coverProtocol(message *entity.Message) (*model.Protocol, *utils.Error) {
	protocol := &model.Protocol{}
	switch message.Type {
	case "friend":
		protocol.Type = 1
		break
	case "group":
		protocol.Type = 8
		break
	}
	protocol.From = message.From
	//从好友中拿he
	friend, e := QueryFriend(message.TargetId, repository.NewFriendRepo())
	if e != nil || friend == nil {
		return nil, log.WithError(utils.ERR_SEND_FAIL_BY_NOT_TARGET)
	}
	protocol.To = strconv.FormatUint(friend.He, 10)
	content, e := util.Obj2Str(message)
	if e != nil {
		return nil, log.WithError(utils.ERR_SEND_FAIL)
	}
	protocol.Data = content
	protocol.Ack = 100
	//生成UUID作为消息号
	protocol.No = message.No
	return protocol, nil
}
func (_self *MessageService) coverMessage(tp string, target uint64, no, content string) *entity.Message {
	message := &entity.Message{}
	message.No = no
	message.Type = tp
	message.TargetId = target
	message.UserId = conf.GetLoginInfo().User.Id
	message.From = strconv.FormatUint(conf.GetLoginInfo().User.Id, 10)
	message.Data = content
	message.Time = uint64(time.Now().Unix())
	message.Ext4 = 1 // 发送中
	message.Ext5 = 2 // 自己发的 肯定是已读
	return message
}
