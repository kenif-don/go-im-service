package service

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/db"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
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

// Paging 消息分页
func (_self *MessageService) Paging(tp string, target, time uint64) ([]entity.Message, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return []entity.Message{}, log.WithError(utils.ERR_NOT_LOGIN)
	}
	pageReq := &entity.Message{
		Type:     tp,
		TargetId: target,
		UserId:   conf.GetLoginInfo().User.Id,
		Time:     time,
	}
	msgs, e := _self.repo.Paging(pageReq)
	if e != nil {
		return []entity.Message{}, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//循环解密
	for i := 0; i < len(msgs); i++ {
		data, err := Decrypt(target, tp, msgs[i].Data)
		if err != nil {
			msgs[i].Data = util.GetErrMsg(utils.ERR_DECRYPT_FAIL)
		} else {
			msgs[i].Data = data
		}
	}
	return msgs, nil
}
func (_self *MessageService) QueryLast(obj *entity.Message) (*entity.Message, error) {
	return _self.repo.QueryLast(obj)
}
func (_self *MessageService) UpdateReaded(protocol *model.Protocol, send int) {
	var message = &entity.Message{}
	e := util.Str2Obj(protocol.Data.(string), message)
	if e != nil {
		log.Error(e)
		return
	}
	//修改消息发送状态
	message.Send = send
	e = _self.repo.Save(message)
	if e != nil {
		log.Error(e)
	}
	//回调
	if Listener != nil {
		resp := &api.SendResp{
			No:   message.No,
			Send: int32(send),
		}
		res, e := proto.Marshal(resp)
		if e != nil {
			log.Error(e)
		}
		Listener.OnSendReceive(res)
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
		if Listener != nil {
			Listener.OnFriendApply()
		}
		break
	case 102: //当to同意好友申请后，更新好友数据
		_, err := NewFriendService().updateOne(util.Str2Uint64(protocol.From), util.Str2Uint64(protocol.To))
		if err != nil {
			return log.WithError(err)
		}
		break
	case 301: // 被好友删除
		err := NewFriendService().DelLocal(&entity.Friend{He: util.Str2Uint64(protocol.From), Me: util.Str2Uint64(protocol.To)})
		return log.WithError(err)
	case 999: //删除聊天和记录
		err := NewChatService().DelChat(protocol.Data.(string), util.Str2Uint64(protocol.From))
		if err != nil {
			return log.WithError(err)
		}
		break
	case 1: // 接收到聊天消息
		//如果是别人发给自己的 就存起来 如果是自己发的 再发送时已经进行了存储
		if util.Str2Uint64(protocol.From) != conf.GetLoginInfo().User.Id {
			messageService := NewMessageService()
			var message = &entity.Message{}
			e := util.Str2Obj(protocol.Data.(string), message)
			if e != nil {
				return log.WithError(e)
			}
			if util.Str2Uint64(protocol.From) == conf.Conf.ChatId {
				//解密
				data, err := Decrypt(util.Str2Uint64(protocol.From), message.Type, message.Data)
				if err != nil {
					return log.WithError(err)
				}
				message.Data = data
				if Listener != nil {
					res, e := util.Obj2Str(message)
					if e != nil {
						return log.WithError(e)
					}
					Listener.OnReceive(res)
				}
			}
			e = messageService.repo.Save(message)
			if e != nil {
				return log.WithError(e)
			}

			//判断是否存在聊天
			chat, e := QueryChat(message.Type, message.UserId, repository.NewChatRepo())
			if e != nil {
				return log.WithError(e)
			}
			if chat == nil {
				chat, e = NewChatService().CoverChat(message.Type, message.UserId)
				if e != nil {
					return log.WithError(e)
				}
			}
			//组装最后一条消息
			err := NewChatService().ChatNotify(chat)
			if err != nil {
				return log.WithError(err)
			}
		}
		break
	}
	return nil
}

func (_self *MessageService) SendMsg(tp string, target uint64, no, content string) *utils.Error {
	//组装本地消息
	message, err := _self.coverMessage(tp, target, no, content)
	if err != nil {
		return log.WithError(err)
	}
	//组装长连接protocol
	protocol, err := _self.coverProtocol(message)
	if err != nil {
		return log.WithError(err)
	}
	//发送消息
	handler.GetClientHandler().GetMessageManager().Send(protocol)
	//保存消息到数据库
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
	protocol.To = strconv.FormatUint(message.TargetId, 10)
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
func (_self *MessageService) coverMessage(tp string, target uint64, no, content string) (*entity.Message, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, log.WithError(utils.ERR_NOT_LOGIN)
	}
	message := &entity.Message{}
	message.No = no
	message.Type = tp
	message.TargetId = target
	message.UserId = conf.GetLoginInfo().User.Id
	message.From = strconv.FormatUint(conf.GetLoginInfo().User.Id, 10)
	//加密
	data, err := Encrypt(message.TargetId, tp, content)
	if err != nil {
		return nil, log.WithError(err)
	}
	message.Data = data
	message.Time = uint64(time.Now().Unix())
	message.Send = 1 // 发送中
	message.Read = 2 // 自己发的 肯定是已读
	return message, nil
}

// Encrypt 聊天内容加密
func Encrypt(he uint64, tp, content string) (string, *utils.Error) {
	key := tp + "_" + util.Uint642Str(he)
	switch tp {
	case "friend":
		if Keys[key] == "" {
			user, e := QueryUser(he, repository.NewUserRepo())
			if e != nil {
				return "", log.WithError(e)
			}
			if user == nil {
				return "", log.WithError(utils.ERR_ENCRYPT_FAIL)
			}
			secret := util.SharedAESKey(user.PublicKey, conf.GetLoginInfo().User.PrivateKey, conf.Conf.Prime)
			Keys[key] = secret
		}
		break
	case "group":
		break
	}
	data, e := util.EncryptAes(content, Keys[key])
	if e != nil {
		return "", log.WithError(e)
	}
	return data, nil
}

// Decrypt 聊天内容解密
func Decrypt(he uint64, tp, content string) (string, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return "", log.WithError(utils.ERR_NOT_LOGIN)
	}
	if content == "" {
		return "", nil
	}
	key := tp + "_" + util.Uint642Str(he)
	switch tp {
	case "friend":
		if Keys[key] == "" {
			user, e := QueryUser(he, repository.NewUserRepo())
			if e != nil {
				return "", log.WithError(e)
			}
			if user == nil {
				return "", log.WithError(utils.ERR_DECRYPT_FAIL)
			}
			secret := util.SharedAESKey(user.PublicKey, conf.GetLoginInfo().User.PrivateKey, conf.Conf.Prime)
			Keys[key] = secret
		}
		break
	case "group":
		break
	}
	data, e := util.DecryptAes(content, Keys[key])
	if e != nil {
		return "", log.WithError(e)
	}
	return data, nil
}
