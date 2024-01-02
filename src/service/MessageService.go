package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/db"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"github.com/google/uuid"
	"im-sdk/model"
	"strconv"
)

type MessageService struct {
	repo *repository.MessageRepo
}

func NewMessageService() *MessageService {
	return &MessageService{
		repo: repository.NewMessageRepo(),
	}
}
func (_self *MessageService) SelectOne(obj *entity.Message) (*entity.Message, *utils.Error) {
	msg, e := _self.repo.Query(obj)
	if e != nil {
		return nil, utils.ERR_MESSAGE_NOT_FOUND
	}
	return msg, nil
}
func (_self *MessageService) Update(obj *entity.Message) *utils.Error {
	e := _self.repo.Save(obj)
	if e != nil {
		return utils.ERR_MESSAGE_UPDATE_FAIL
	}
	return nil
}

// DelChatMsg 删除双方聊天记录
func (_self *MessageService) DelChatMsg(tp string, target uint64) *utils.Error {
	//组装data数据
	data := make(map[string]string)
	switch tp {
	case "friend":
		data["target"] = util.Uint642Str(conf.GetLoginInfo().User.Id)
		break
	case "group":
		data["target"] = util.Uint642Str(target)
		break
	}
	data["type"] = tp
	dataStr, e := util.Obj2Str(data)
	if e != nil {
		log.Error(e)
		return log.WithError(utils.ERR_DEL_FAIL)
	}
	//发出消息 让对方删除
	//发送删除请求
	protocol := &model.Protocol{
		Type: 998,
		From: util.Uint642Str(conf.GetLoginInfo().User.Id),
		To:   util.Uint642Str(target),
		Ack:  100,
		Data: dataStr,
		No:   uuid.New().String(),
	}
	err := Send(protocol)
	if err != nil {
		return log.WithError(err)
	}
	return _self.DelLocalChatMsg(tp, target)
}

// DelLocalChatMsg 删除登录者本地消息
func (_self *MessageService) DelLocalChatMsg(tp string, target uint64) *utils.Error {
	tx := _self.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		//自己删除自己发送给对方的 或对方发给自己的
		message := &entity.Message{
			Type:     tp,
			TargetId: target,
			UserId:   conf.GetLoginInfo().User.Id,
		}
		e := _self.repo.Delete(message)
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		//如果是PC 更新会话
		if conf.Base.DeviceType == conf.PC {
			//更新会话
			err := NewChatService().ChatNotify(&entity.Chat{
				Type:     tp,
				TargetId: target,
				UserId:   conf.GetLoginInfo().User.Id,
			})
			if err != nil {
				return log.WithError(utils.ERR_DEL_FAIL)
			}
		}
		//如果当前打开的会话是要被删除聊天记录的 就进行通知
		err := DelMsgNotify(tp, target)
		if err != nil {
			return log.WithError(err)
		}
		e = tx.Commit().Error
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// DelAllMessage 删除制定用户所有消息 用于自毁
func (_self *MessageService) DelAllMessage() *utils.Error {
	message := &entity.Message{
		UserId: conf.GetLoginInfo().User.Id,
	}
	e := _self.repo.Delete(message)
	if e != nil {
		return log.WithError(utils.ERR_DEL_FAIL)
	}
	return nil
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
		log.Error(e)
		return []entity.Message{}, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//循环解密
	for i := 0; i < len(msgs); i++ {
		data, err := Decrypt(tp, target, msgs[i].No, msgs[i].Data)
		if err != nil {
			msgs[i].Data = util.GetTextErrMsg()
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
}

// GetOfflineMessage 获取离线消息
func (_self *MessageService) GetOfflineMessage() *utils.Error {
	resultDTO, err := Post("/api/offline-bill/selectAll", nil)
	if err != nil {
		return log.WithError(err)
	}
	var offlineBills = &[]entity.OfflineBill{}
	e := util.Str2Obj(resultDTO.Data.(string), offlineBills)
	if e != nil {
		return log.WithError(utils.ERR_QUERY_FAIL)
	}
	tx := db.NewTransaction().BeginTx()
	if e := tx.Error; e != nil {
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
		_, err = Post("/api/offline-bill/dels", map[string]interface{}{"ids": ids})
		if err != nil {
			return log.WithError(err)
		}
		e = tx.Commit().Error
		if e != nil {
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}
func (_self *MessageService) Handler(protocol *model.Protocol) *utils.Error {
	switch protocol.Type {
	case 101: //让to拉去from的好友申请信息，没有就存起来 有就修改
		friendAppluService := NewFriendApplyService()
		err := friendAppluService.updateOne(util.Str2Uint64(protocol.From), util.Str2Uint64(protocol.To))
		if err != nil {
			return log.WithError(err)
		}
		//发送通知
		friendAppluService.FriendApplyNotify()
		break
	case 102: //当to同意好友申请后，更新好友数据
		_, err := NewFriendService().SelectOne(util.Str2Uint64(protocol.From), false)
		if err != nil {
			return log.WithError(err)
		}
		break
	case 201: // 系统指令 去服务器拉去群成员
		gId := util.Str2Uint64(protocol.Data.(string))
		_, err := NewGroupMemberService().SelectMembers(gId, true)
		if err != nil {
			return log.WithError(err)
		}
		break
	case 209: // 删除本地群
		gId := util.Str2Uint64(protocol.Data.(string))
		err := NewGroupService().DelLocalGroup(gId)
		if err != nil {
			return log.WithError(err)
		}
		break
	case 301: // 被好友删除
		err := NewFriendService().DelLocalFriend(&entity.Friend{He: util.Str2Uint64(protocol.From), Me: util.Str2Uint64(protocol.To)})
		if err != nil {
			return log.WithError(err)
		}
		break
	case 403: // 客户端被T 直接退出登录
		log.Debug("客户端被T 直接退出登录")
		err := NewUserService().Logout()
		if err != nil {
			return log.WithError(err)
		}
	case 998: //被删除本地聊天记录
		data := make(map[string]string)
		log.Debugf("被删除本地聊天记录%v", protocol.Data)
		e := util.Str2Obj(protocol.Data.(string), &data)
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		err := NewMessageService().DelLocalChatMsg(data["type"], util.Str2Uint64(data["target"]))
		if err != nil {
			return log.WithError(err)
		}
		break
	case 999: //被删除聊天和记录
		data := make(map[string]string)
		log.Debugf("被删除聊天和记录%v", protocol.Data)
		e := util.Str2Obj(protocol.Data.(string), &data)
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		err := NewChatService().DelLocalChat(data["type"], util.Str2Uint64(data["target"]))
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
			//重置userId为当前用户 不然userId就是发送者了
			message.UserId = conf.GetLoginInfo().User.Id
			//别人给自己发的 肯定是发送成功
			message.Send = 2
			e = messageService.repo.Save(message)
			if e != nil {
				return log.WithError(e)
			}
			//如果发送者是当前用户打开的聊天目标
			if util.Str2Uint64(protocol.From) == conf.Conf.ChatId {
				//解密
				data, err := Decrypt(message.Type, util.Str2Uint64(protocol.From), message.No, message.Data)
				if err != nil {
					data = util.GetTextErrMsg()
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
			//判断是否存在聊天
			chat, e := QueryChat(message.Type, util.Str2Uint64(protocol.From), repository.NewChatRepo())
			if e != nil {
				return log.WithError(e)
			}
			if chat == nil {
				c, err := NewChatService().CoverChat(message.Type, util.Str2Uint64(protocol.From), false)
				if err != nil {
					return log.WithError(err)
				}
				chat = c
			}
			// 通知聊天列表更新
			err := NewChatService().ChatNotify(chat)
			if err != nil {
				return log.WithError(err)
			}
			// 通知语音播报
			err = NewChatService().VoiceNotify(message)
		}
		break
	case 8: //群聊消息
		if util.Str2Uint64(protocol.From) != conf.GetLoginInfo().User.Id {
			messageService := NewMessageService()
			var message = &entity.Message{}
			e := util.Str2Obj(protocol.Data.(string), message)
			if e != nil {
				return log.WithError(e)
			}
			message.No = protocol.No
			//重置userId为当前用户 不然userId就是发送者了
			message.UserId = conf.GetLoginInfo().User.Id
			//别人给自己发的 肯定是发送成功
			message.Send = 2
			e = messageService.repo.Save(message)
			if e != nil {
				return log.WithError(e)
			}
			//判断是否存在聊天
			chat, e := QueryChat(message.Type, message.TargetId, repository.NewChatRepo())
			if e != nil {
				return log.WithError(e)
			}
			if chat == nil {
				c, err := NewChatService().CoverChat(message.Type, message.TargetId, false)
				if err != nil {
					return log.WithError(err)
				}
				chat = c
			}
			// 通知聊天列表更新
			err := NewChatService().ChatNotify(chat)
			if err != nil {
				return log.WithError(err)
			}
			// 通知语音播报
			err = NewChatService().VoiceNotify(message)
			//如果发送者是当前用户打开的聊天目标
			if message.TargetId == conf.Conf.ChatId {
				//解密
				data, err := Decrypt(message.Type, message.TargetId, message.No, message.Data)
				if err != nil {
					data = util.GetTextErrMsg()
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
		}
		break
	}
	return nil
}

func (_self *MessageService) SendMsg(tp string, target uint64, no string, dataContent *entity.MessageData) *utils.Error {
	//判断好友或者群是否存在
	switch tp {
	case "friend":
		//先本地查
		friend, err := NewFriendService().SelectOne(target, false)
		if err != nil || friend == nil {
			return log.WithError(utils.ERR_FRIEND_GET_FAIL)
		}
		break
	case "group":
		//先本地查
		friend, err := NewGroupService().SelectOne(target, false)
		if err != nil || friend == nil {
			return log.WithError(utils.ERR_FRIEND_GET_FAIL)
		}
	}
	switch dataContent.Type {
	case 1: //文本消息
		return _self.realSend(tp, target, no, dataContent)
	case 2, 5: //图片消息/文件消息
		return _self.SendImgAndFileMsg(tp, target, no, dataContent)
	case 3: //语音消息
		return _self.SendVoiceMsg(tp, target, no, dataContent)
	case 4: //视频消息
		return _self.SendVideoMsg(tp, target, no, dataContent)
	}
	return nil
}
func (_self *MessageService) SendVideoMsg(tp string, target uint64, no string, dataContent *entity.MessageData) *utils.Error {
	secret, err := GetSecret(target, tp)
	if err != nil {
		return log.WithError(err)
	}
	//上传文件
	url, err := util.Upload(dataContent.Content, secret)
	if err != nil {
		return log.WithError(err)
	}
	dataContent.Content = url
	return _self.realSend(tp, target, no, dataContent)
}
func (_self *MessageService) SendVoiceMsg(tp string, target uint64, no string, dataContent *entity.MessageData) *utils.Error {
	secret, err := GetSecret(target, tp)
	if err != nil {
		return log.WithError(err)
	}
	//上传文件
	url, err := util.Upload(dataContent.Content, secret)
	if err != nil {
		return log.WithError(err)
	}
	dataContent.Content = url
	return _self.realSend(tp, target, no, dataContent)
}
func (_self *MessageService) SendImgAndFileMsg(tp string, target uint64, no string, dataContent *entity.MessageData) *utils.Error {
	secret, err := GetSecret(target, tp)
	if err != nil {
		return log.WithError(err)
	}
	//上传文件
	url, err := util.Upload(dataContent.Content, secret)
	if err != nil {
		return log.WithError(err)
	}
	dataContent.Content = url
	return _self.realSend(tp, target, no, dataContent)
}

// realSend 发送文本消息
func (_self *MessageService) realSend(tp string, target uint64, no string, dataContent *entity.MessageData) *utils.Error {
	//组装本地消息
	message, err := _self.coverMessage(tp, target, no, dataContent)
	if err != nil {
		return log.WithError(err)
	}
	//组装长连接protocol
	protocol, err := _self.coverProtocol(message)
	if err != nil {
		return log.WithError(err)
	}
	//先保存消息到数据库
	e := _self.repo.Save(message)
	if e != nil {
		return log.WithError(utils.ERR_SEND_FAIL)
	}
	//再发送消息
	err = Send(protocol)
	if err != nil {
		return log.WithError(err)
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
func (_self *MessageService) coverMessage(tp string, target uint64, no string, dataContent *entity.MessageData) (*entity.Message, *utils.Error) {
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
	data, err := Encrypt(message.TargetId, tp, dataContent.Content)
	if err != nil {
		log.Error(err)
		return nil, log.WithError(utils.ERR_SEND_FAIL)
	}
	dataContent.Content = data
	encryptData, e := util.Obj2Str(dataContent)
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_SEND_FAIL)
	}
	message.Data = encryptData
	message.Time = _self.CurrentTime()
	message.Send = 1 // 1-发送中 2-发送成功 -1-发送失败
	message.Read = 2 // 自己发的 肯定是已读
	return message, nil
}
func (_self *MessageService) CurrentTime() uint64 {
	return uint64(int64(util.CurrentTime()) + int64(conf.DiffTime))
}
