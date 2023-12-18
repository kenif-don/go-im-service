package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/db"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"gorm.io/gorm"
)

type IFriendRepo interface {
	Query(obj *entity.Friend) (*entity.Friend, error)
	QueryAll(obj *entity.Friend) ([]entity.Friend, error)
	Save(obj *entity.Friend) error
	Delete(obj *entity.Friend) error
	BeginTx() *gorm.DB
}
type FriendService struct {
	repo IFriendRepo
}

func NewFriendService() *FriendService {
	return &FriendService{
		repo: repository.NewFriendRepo(),
	}
}

// SelectOne 获取好友 逻辑： 先从本地获取 获取失败或需要刷新 则从远程获取
func (_self *FriendService) SelectOne(he uint64, refresh bool) (*entity.Friend, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, log.WithError(utils.ERR_NOT_LOGIN)
	}
	me := conf.GetLoginInfo().User.Id
	//先从本地获取
	friend, e := _self.repo.Query(&entity.Friend{He: he, Me: me})
	if e != nil {
		return nil, log.WithError(utils.ERR_FRIEND_GET_FAIL)
	}
	//没有获取到 或者需要刷新好友数据
	if friend == nil || refresh {
		resultDTO, err := Post("/api/friend/selectOne", map[string]uint64{"from": he, "to": me})
		if err != nil {
			return nil, log.WithError(err)
		}
		//如果服务器获取失败
		if resultDTO.Data == nil {
			return nil, log.WithError(utils.ERR_FRIEND_GET_FAIL)
		}
		var fa entity.Friend
		e := util.Str2Obj(resultDTO.Data.(string), &fa)
		if e != nil {
			return nil, log.WithError(utils.ERR_FRIEND_GET_FAIL)
		}
		if fa.Id != 0 {
			//保存到数据库
			e := _self.repo.Save(&fa)
			if e != nil {
				return nil, log.WithError(utils.ERR_OPERATION_FAIL)
			}
			if fa.HeUser == nil {
				return nil, nil
			}
			//再保存好友用户
			user, err := NewUserService().SelectOne(fa.He, refresh)
			if err != nil {
				return nil, log.WithError(utils.ERR_OPERATION_FAIL)
			}
			fa.HeUser = user
		}
		return &fa, nil
	}
	//获取到了 组装heUser
	if friend != nil {
		user, err := NewUserService().SelectOne(friend.He, false)
		if err != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		friend.HeUser = user
	}
	return friend, nil
}
func (_self *FriendService) SelectAll() ([]entity.Friend, *utils.Error) {
	friends, e := _self.repo.QueryAll(&entity.Friend{Me: conf.GetLoginInfo().User.Id})
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	if len(friends) != 0 {
		//封装好友信息
		for i := 0; i < len(friends); i++ {
			user, err := NewUserService().SelectOne(friends[i].He, false)
			if err != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			friends[i].HeUser = user
			if friends[i].Name == "" {
				friends[i].Name = user.Nickname
			}
		}
		return friends, nil
	}
	//没查到 就从后台查一次
	resultDTO, err := Post("/api/friend/selectAll", nil)
	if err != nil {
		return nil, log.WithError(err)
	}
	var fs []entity.Friend
	e = util.Str2Obj(resultDTO.Data.(string), &fs)
	if e != nil || fs == nil {
		return []entity.Friend{}, nil
	}
	tx := _self.repo.BeginTx()
	if e := tx.Error; e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	r, err := func() ([]entity.Friend, *utils.Error) {
		// 保存到数据库
		for _, v := range fs {
			e := _self.repo.Save(&v)
			if e != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			//保存对应的用户信息
			_, err := NewUserService().SelectOne(v.He, false)
			if err != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
		}
		e = tx.Commit().Error
		if e != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		return fs, nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return r, err
}

// DelFriend 删除双方好友
func (_self *FriendService) DelFriend(id uint64) *utils.Error {
	//先删本地
	err := _self.DelLocalFriend(&entity.Friend{He: id, Me: conf.GetLoginInfo().User.Id})
	if err != nil {
		return log.WithError(err)
	}
	//再通过服务器删除 这里服务器删除的就是双方的 服务器回去发送长连接
	req := make(map[string]uint64)
	req["he"] = id
	_, err = Post("/api/friend/delete", req)
	if err != nil {
		return log.WithError(err)
	}
	return nil
}
func (_self *FriendService) DelLocalFriend(friend *entity.Friend) *utils.Error {
	tx := db.NewTransaction().BeginTx()
	if e := tx.Error; e != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		f, e := _self.repo.Query(friend)
		if e != nil {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		//先修改好友申请为拒绝
		faService := NewFriendApplyService()
		err := faService.updateReject(f.Me, f.He)
		if err != nil {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		err = faService.updateReject(f.He, f.Me)
		if err != nil {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		//删除聊天--删除聊天时即会删除消息,也会对PC聊天列表进行通知
		err = NewChatService().DelLocalChat("friend", f.He)
		if err != nil {
			return log.WithError(err)
		}
		//再删除本地好友记录 不删用户
		e = _self.repo.Delete(f)
		if e != nil {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		e = tx.Commit().Error
		if e != nil {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}

func (_self *FriendService) UpdateName(he uint64, name string) *utils.Error {
	friend, err := _self.SelectOne(he, false)
	if err != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	friend.Name = name
	//服务器修改
	_, err = Post("/api/friend/edit", friend)
	if err != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	e := _self.repo.Save(friend)
	if e != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	//如果聊天存在 同步修改聊天中的name
	chat, e := QueryChat("friend", friend.He, repository.NewChatRepo())
	if e != nil {
		log.Error(e)
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	if chat != nil {
		chat.Name = name
		e = NewChatService().repo.Save(chat)
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		//通知客户端更新聊天列表
		if conf.Base.DeviceType == conf.PC {
			err = NewChatService().ChatNotify(&entity.Chat{
				Type:     "friend",
				TargetId: friend.He,
			})
			if err != nil {
				return log.WithError(utils.ERR_OPERATION_FAIL)
			}
		}
	}
	return nil
}
