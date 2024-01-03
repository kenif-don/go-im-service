package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"im-sdk/model"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	once sync.Once
	Keys map[string]string // 加密缓存 对方用户id为key，秘钥为value
)

type IChatRepo interface {
	Query(obj *entity.Chat) (*entity.Chat, error)
	QueryAll(userId uint64) ([]entity.Chat, error)
	Save(obj *entity.Chat) error
	Delete(obj *entity.Chat) error
	BeginTx() *gorm.DB
}
type ChatService struct {
	repo IChatRepo
}

func NewChatService() *ChatService {
	return &ChatService{
		repo: repository.NewChatRepo(),
	}
}
func QueryChat(tp string, target uint64, repo IChatRepo) (*entity.Chat, error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, log.WithError(utils.ERR_NOT_LOGIN)
	}
	return repo.Query(&entity.Chat{Type: tp, TargetId: target, UserId: conf.GetLoginInfo().User.Id})
}
func (_self *ChatService) OpenChat(tp string, target uint64, password string) (*entity.Chat, *utils.Error) {
	//再根据最新user 更新一次聊天信息 如果用户更新了昵称 头像等 这里可以刷新
	chat, err := _self.CoverChat(tp, target, true)
	if err != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//根据类型查询数据
	switch tp {
	case "friend":
		//清除秘钥缓存
		Keys["friend"+"_"+util.Uint642Str(target)] = ""
		break
	case "group":
		//如果是加密群
		group, e := NewGroupService().SelectOne(target, false)
		if e != nil {
			log.Error(e)
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		if group.Type != 2 {
			break
		}
		//没有输入过 并且没有传进来 就提示需要输入密码
		if conf.Conf.Pwds[tp+"_"+util.Uint642Str(target)] == "" && password == "" {
			return nil, log.WithError(utils.ERR_ENTER_PASSWORD)
			//用户传了密码 但是和数据库不一致
		} else if password != "" && strings.ToUpper(group.Password) != strings.ToUpper(util.MD5(password)) {
			return nil, log.WithError(utils.ERR_PASSWORD_ERROR)
			//用户传了密码 但是与内存不一致
		} else if password != "" && conf.Conf.Pwds[tp+"_"+util.Uint642Str(target)] != "" && conf.Conf.Pwds[tp+"_"+util.Uint642Str(target)] != password {
			return nil, log.WithError(utils.ERR_PASSWORD_ERROR)
		} else if password != "" {
			conf.Conf.Pwds[tp+"_"+util.Uint642Str(target)] = password
		}
		break
	}
	//记录当前聊天ID
	conf.Conf.ChatId = chat.TargetId
	//如果是PC的话 需要通知客户端更新聊天列表
	if conf.Base.DeviceType == conf.PC {
		err := _self.ChatNotify(chat)
		if err != nil {
			return nil, log.WithError(err)
		}
	}
	return chat, nil
}

// CoverChat 封装聊天
func (_self *ChatService) CoverChat(tp string, target uint64, refresh bool) (*entity.Chat, *utils.Error) {
	var name, headImg string
	switch tp {
	case "friend":
		//获取好友信息
		friend, err := NewFriendService().SelectOne(target, refresh)
		if err != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		//组装聊天名称
		if friend.Name != "" {
			name = friend.Name
		} else {
			name = friend.HeUser.Nickname
		}
		headImg = friend.HeUser.HeadImg
		break
	case "group":
		//获取群信息
		group, err := NewGroupService().SelectOne(target, refresh)
		if err != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		name = group.Name
		headImg = group.HeadImg
		break
	}
	chat := &entity.Chat{
		Type:     tp,
		TargetId: target,
		UserId:   conf.GetLoginInfo().User.Id,
		Name:     name,
		HeadImg:  headImg,
		UnReadNo: 0,
	}
	e := _self.repo.Save(chat)
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	return chat, nil
}
func (_self *ChatService) GetChats() (*[]entity.Chat, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, log.WithError(utils.ERR_NOT_LOGIN)
	}
	chats, err := _self.repo.QueryAll(conf.GetLoginInfo().User.Id)
	if err != nil {
		return &[]entity.Chat{}, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//组装所有消息和最后一条消息
	for i := range chats {
		err := _self.coverLastMsg(&chats[i])
		if err != nil {
			return &[]entity.Chat{}, log.WithError(utils.ERR_QUERY_FAIL)
		}
	}
	//排序 根据最后的消息时间倒序
	sort.Slice(chats, func(i, j int) bool {
		return chats[j].LastTime < chats[i].LastTime
	})
	return &chats, nil
}
func (_self *ChatService) coverLastMsg(chat *entity.Chat) *utils.Error {
	//组装最后一条消息
	messageService := NewMessageService()
	pageReq := &entity.Message{
		Type:     chat.Type,
		TargetId: chat.TargetId,
		UserId:   conf.GetLoginInfo().User.Id,
	}
	lastMsg, e := messageService.QueryLast(pageReq)
	if e != nil {
		return log.WithError(utils.ERR_QUERY_FAIL)
	}
	if lastMsg != nil {
		//解密
		data, err := Decrypt(chat.Type, chat.TargetId, "", lastMsg.Data)
		if err != nil {
			chat.LastMsg = util.GetTextErrMsg()
		} else {
			chat.LastMsg = data
		}
		chat.LastTime = lastMsg.Time
	}
	return nil
}

// DelLocalChat 删除本地聊天记录
func (_self *ChatService) DelLocalChat(tp string, target uint64) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return log.WithError(utils.ERR_NOT_LOGIN)
	}
	tx := _self.repo.BeginTx()
	if e := tx.Error; e != nil {
		return log.WithError(utils.ERR_DEL_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		log.Debugf("del local chat %s %d", tp, target)
		// 删除聊天
		var chat entity.Chat
		chat.Type = tp
		chat.TargetId = target
		chat.UserId = conf.GetLoginInfo().User.Id
		e := _self.repo.Delete(&chat)
		if e != nil {
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		log.Debugf("del local message %s %d", tp, target)
		// 删除聊天记录
		err := NewMessageService().DelLocalChatMsg(tp, target)
		if err != nil {
			log.Error(err)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		e = tx.Commit().Error
		if e != nil {
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		log.Debugf("提交事务111111")
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// DelChat 删除双方聊天记录
func (_self *ChatService) DelChat(tp string, target uint64) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return log.WithError(utils.ERR_NOT_LOGIN)
	}
	//发送删除请求
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
	protocol := &model.Protocol{
		Type: 999,
		From: strconv.FormatUint(conf.GetLoginInfo().User.Id, 10),
		To:   strconv.FormatUint(target, 10),
		Ack:  100,
		Data: dataStr,
		No:   uuid.New().String(),
	}
	err := Send(protocol)
	if err != nil {
		return log.WithError(err)
	}
	return _self.DelLocalChat(tp, target)
}

// DelAllChat 删除制定用户所有聊天
func (_self *ChatService) DelAllChat() *utils.Error {
	chat := &entity.Chat{
		UserId: conf.GetLoginInfo().User.Id,
	}
	e := _self.repo.Delete(chat)
	if e != nil {
		return log.WithError(utils.ERR_DEL_FAIL)
	}
	return nil
}
