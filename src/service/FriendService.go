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
func QueryFriend(obj *entity.Friend, repo IFriendRepo) (*entity.Friend, error) {
	return repo.Query(obj)
}
func QueryFriendAll(repo IFriendRepo) ([]entity.Friend, error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return []entity.Friend{}, nil
	}
	return repo.QueryAll(&entity.Friend{Me: conf.GetLoginInfo().User.Id})
}

func (_self *FriendService) updateOne(he, me uint64) (*entity.Friend, *utils.Error) {
	var req = make(map[string]uint64)
	req["he"] = he
	req["me"] = me
	resultDTO, err := Post("/api/friend/selectOne", req)
	if err != nil {
		return nil, log.WithError(err)
	}
	var fa entity.Friend
	_ = util.Str2Obj(resultDTO.Data.(string), &fa)
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
		userService := NewUserService()
		e = userService.Save(fa.HeUser)
		if e != nil {
			return nil, log.WithError(utils.ERR_OPERATION_FAIL)
		}
	}
	return &fa, nil
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
func (_self *FriendService) QueryFriend2(he uint64) (*entity.Friend, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, log.WithError(utils.ERR_NOT_LOGIN)
	}
	friend, e := _self.repo.Query(&entity.Friend{He: he, Me: conf.GetLoginInfo().User.Id})
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	if friend != nil {
		userService := NewUserService()
		user, e := QueryUser(friend.He, userService.repo)
		if e != nil || user == nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		friend.HeUser = user
	}
	return friend, nil
}
func (_self *FriendService) SelectOne(he uint64) (*entity.Friend, *utils.Error) {
	friend, e := QueryFriend(&entity.Friend{He: he, Me: conf.GetLoginInfo().User.Id}, _self.repo)
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	if friend != nil {
		userService := NewUserService()
		user, e := QueryUser(friend.He, userService.repo)
		if e != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		if user == nil {
			f, err := _self.updateOne(friend.He, friend.Me)
			if err != nil {
				return nil, log.WithError(err)
			}
			user = f.HeUser
		}
		friend.HeUser = user
	}
	return friend, nil
}
func (_self *FriendService) SelectAll() ([]entity.Friend, *utils.Error) {
	friends, e := QueryFriendAll(_self.repo)
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	if len(friends) != 0 {
		//封装好友信息
		for i := 0; i < len(friends); i++ {
			userService := NewUserService()
			user, e := QueryUser(friends[i].He, userService.repo)
			if e != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			if user == nil {
				f, err := _self.updateOne(friends[i].He, friends[i].Me)
				if err != nil {
					return nil, log.WithError(err)
				}
				user = f.HeUser
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
			//先查询 是否存在 存在就不添加了
			//保存对应的用户信息
			userService := NewUserService()
			sysUser, e := QueryUser(v.He, userService.repo)
			if e != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			if sysUser != nil {
				continue
			}
			e = userService.Save(v.HeUser)
			if e != nil {
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
func (_self *FriendService) UpdateName(id uint64, name string) *utils.Error {
	friend, e := QueryFriend(&entity.Friend{He: id, Me: conf.GetLoginInfo().User.Id}, _self.repo)
	if e != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	friend.Name = name
	//服务器修改
	_, err := Post("/api/friend/edit", friend)
	if err != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	e = _self.repo.Save(friend)
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
