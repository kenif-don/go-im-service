package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"gorm.io/gorm"
	"sync"
)

var (
	once     sync.Once
	Listener MessageListener
)

type MessageListener interface {
	//OnReceive 当前聊天接收到消息
	OnReceive(data []byte)
	//OnSendReceive 发送的消息状态 -某消息 发送成功、发送失败
	OnSendReceive(data []byte)
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

func NewChatService(listener MessageListener) *ChatService {
	once.Do(func() {
		Listener = listener
	})
	return &ChatService{
		repo: repository.NewChatRepo(),
	}
}
func QueryChatOne(id uint64, repo IChatRepo) (*entity.Chat, error) {
	return repo.Query(&entity.Chat{Id: id})
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
		return nil, utils.ERR_QUERY_FAIL
	}
	if chat == nil {
		//根据类型查询数据
		switch tp {
		case "friend":
			friend, err := NewFriendService().SelectOne(target)
			if err != nil || friend == nil {
				return nil, utils.ERR_QUERY_FAIL
			}
			var name string
			if friend.Name != "" {
				name = friend.Name
			} else {
				name = friend.HeUser.Nickname
			}
			chat = &entity.Chat{
				Type:     tp,
				TargetId: friend.Id,
				UserId:   conf.GetLoginInfo().User.Id,
				Name:     name,
				HeadImg:  friend.HeUser.HeadImg,
				UnReadNo: 0,
			}
			e := _self.repo.Save(chat)
			if e != nil {
				return nil, utils.ERR_QUERY_FAIL
			}
			break
		case "group":
			break
		}
	}
	//组装最后一条消息和最新15条消息
	messageService := NewMessageService()
	pageReq := &entity.Message{
		Type:     tp,
		TargetId: target,
		UserId:   conf.GetLoginInfo().User.Id,
	}
	lastMsg, e := messageService.QueryLast(pageReq)
	if e != nil {
		return nil, utils.ERR_QUERY_FAIL
	}
	if lastMsg != nil {
		chat.LastMsg = lastMsg.Data
		chat.LastTime = lastMsg.Time
	}
	chat.Page = 1
	chat.TotalPage = messageService.repo.CountPage(pageReq)
	msgs, e := messageService.repo.Paging(pageReq, chat.Page)
	if e != nil {
		return nil, utils.ERR_QUERY_FAIL
	}
	chat.Msgs = msgs
	//记录当前聊天ID
	conf.Conf.ChatId = chat.Id
	return chat, nil
}
