package service

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"im-sdk/handler"
	"im-sdk/model"
	"sort"
	"strconv"
	"sync"
)

var (
	once     sync.Once
	Listener MessageListener
	Keys     map[string]string // 加密缓存 对方用户id为key，秘钥为value
)

type MessageListener interface {
	//OnReceive 当前聊天接收到消息
	OnReceive(data []byte)
	//OnSendReceive 发送的消息状态 -某消息 发送成功、发送失败
	OnSendReceive(data []byte)
	//OnDoChats 如果客户端停留在首页 如果有新消息进来,都会调用此接口更新最后消息和排序
	OnDoChats(data []byte)
}

func SetListener(listener MessageListener) {
	once.Do(func() {
		Listener = listener
		Keys = make(map[string]string)
	})
}

type IChatRepo interface {
	Query(obj *entity.Chat) (*entity.Chat, error)
	QueryAll(obj *entity.Chat) ([]entity.Chat, error)
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
		return nil, nil
	}
	return repo.Query(&entity.Chat{Type: tp, TargetId: target, UserId: conf.GetLoginInfo().User.Id})
}
func (_self *ChatService) OpenChat(tp string, target uint64) (*entity.Chat, *utils.Error) {
	chat, e := QueryChat(tp, target, _self.repo)
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	if chat == nil {
		//根据类型查询数据
		switch tp {
		case "friend":
			friend, err := NewFriendService().QueryFriend2(target)
			if err != nil || friend == nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			var name string
			if friend.Name != "" {
				name = friend.Name
			} else {
				name = friend.HeUser.Nickname
			}
			chat = &entity.Chat{
				Type:     tp,
				TargetId: target,
				UserId:   conf.GetLoginInfo().User.Id,
				Name:     name,
				HeadImg:  friend.HeUser.HeadImg,
				UnReadNo: 0,
			}
			e := _self.repo.Save(chat)
			if e != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			break
		case "group":
			break
		}
	}
	err := _self.coverMsgs(chat)
	if err != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	err = _self.coverLastMsg(chat)
	if err != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//记录当前聊天ID
	conf.Conf.ChatId = chat.TargetId
	return chat, nil
}

func (_self *ChatService) GetChats() (*[]entity.Chat, *utils.Error) {
	chats, err := _self.repo.QueryAll(&entity.Chat{
		UserId: conf.GetLoginInfo().User.Id,
	})
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
		return chats[i].LastTime < chats[j].LastTime
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
	log.Debugf("最后一条消息:%v", lastMsg)
	if lastMsg != nil {
		//解密
		data, err := Decrypt(chat.TargetId, chat.Type, lastMsg.Data)
		if err != nil {
			chat.LastMsg = util.GetErrMsg(utils.ERR_DECRYPT_FAIL)
		} else {
			chat.LastMsg = data
		}
		chat.LastTime = lastMsg.Time
	}
	return nil
}
func (_self *ChatService) coverMsgs(chat *entity.Chat) *utils.Error {
	//最新15条消息
	messageService := NewMessageService()
	pageReq := &entity.Message{
		Type:     chat.Type,
		TargetId: chat.TargetId,
		UserId:   conf.GetLoginInfo().User.Id,
	}
	chat.Page = 1
	chat.TotalPage = messageService.repo.CountPage(pageReq)
	msgs, e := messageService.repo.Paging(pageReq, chat.Page)
	if e != nil {
		return log.WithError(utils.ERR_QUERY_FAIL)
	}
	//循环解密
	for i := 0; i < len(msgs); i++ {
		data, err := Decrypt(chat.TargetId, chat.Type, msgs[i].Data)
		if err != nil {
			msgs[i].Data = util.GetErrMsg(utils.ERR_DECRYPT_FAIL)
		} else {
			msgs[i].Data = data
		}
	}
	chat.Msgs = msgs
	return nil
}

// DelLocalChat 删除本地聊天记录
func (_self *ChatService) DelLocalChat(tp string, target uint64) *utils.Error {
	tx := _self.repo.BeginTx()
	if err := tx.Error; err != nil {
		return log.WithError(utils.ERR_DEL_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		// 删除聊天
		var chat entity.Chat
		chat.Type = tp
		chat.TargetId = target
		chat.UserId = conf.GetLoginInfo().User.Id
		e := _self.repo.Delete(&chat)
		if e != nil {
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		// 删除聊天记录
		var message entity.Message
		message.Type = tp
		message.TargetId = target
		message.UserId = conf.GetLoginInfo().User.Id
		e = NewMessageService().repo.Delete(&message)
		if e != nil {
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		//调用事件通知
		err := NewChatService().ChatsNotify()
		if err != nil {
			return log.WithError(err)
		}
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return nil
}

// DelChat 删除双方聊天记录
func (_self *ChatService) DelChat(tp string, target uint64) *utils.Error {
	//发送删除请求
	protocol := &model.Protocol{
		Type: 999,
		From: strconv.FormatUint(conf.GetLoginInfo().User.Id, 10),
		To:   strconv.FormatUint(target, 10),
		Ack:  100,
		Data: tp, //将聊天类型传递过去
		No:   uuid.New().String(),
	}
	handler.GetClientHandler().GetMessageManager().Send(protocol)
	err := _self.DelLocalChat(tp, target)
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

// ChatsNotify 通知客户端更新聊天列表
func (_self *ChatService) ChatsNotify() *utils.Error {
	chats, err := _self.GetChats()
	if err != nil {
		return log.WithError(err)
	}
	if Listener != nil {
		resp := &api.ChatData{}
		e := util.Obj2Obj(chats, resp)
		if e != nil {
			return log.WithError(e)
		}
		res, e := proto.Marshal(resp)
		if e != nil {
			return log.WithError(e)
		}
		Listener.OnDoChats(res)
	}
	return nil
}
