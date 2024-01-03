package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"gorm.io/gorm"
)

type IGroupRepo interface {
	Query(obj *entity.Group) (*entity.Group, error)
	QueryAll(obj *entity.Group) ([]entity.Group, error)
	Save(obj *entity.Group) error
	Delete(obj *entity.Group) error
	BeginTx() *gorm.DB
}
type GroupService struct {
	repo IGroupRepo
}

func NewGroupService() *GroupService {
	return &GroupService{
		repo: repository.NewGroupRepo(),
	}
}
func QueryGroup(obj *entity.Group, repo IGroupRepo) (*entity.Group, error) {
	return repo.Query(obj)
}

// Invite 邀请好友进群 ids是用户ID
func (_self *GroupService) Invite(id uint64, ids []uint64) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return log.WithError(utils.ERR_NOT_LOGIN)
	}
	//先从服务器创建
	_, err := Post("/api/group/invite", map[string]interface{}{"ids": ids, "id": id})
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

// Create 创建群聊
func (_self *GroupService) Create(ids []uint64, tp int, password string) (*entity.Group, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, log.WithError(utils.ERR_NOT_LOGIN)
	}
	//先从服务器创建
	resultDTO, err := Post("/api/group/create", map[string]interface{}{"ids": ids, "type": tp, "password": password})
	if err != nil {
		return nil, log.WithError(err)
	}
	var group entity.Group
	e := util.Str2Obj(resultDTO.Data.(string), &group)
	if e != nil {
		return nil, log.WithError(utils.ERR_OPERATION_FAIL)
	}
	if group.Id != 0 {
		//保存到数据库
		e := _self.repo.Save(&group)
		if e != nil {
			return nil, log.WithError(utils.ERR_OPERATION_FAIL)
		}
	}
	return &group, nil
}
func (_self *GroupService) SelectAll() ([]entity.Group, error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return []entity.Group{}, nil
	}
	return _self.repo.QueryAll(&entity.Group{UserId: conf.GetLoginInfo().User.Id})
}

func (_self *GroupService) SelectOne(target uint64, refresh bool) (*entity.Group, *utils.Error) {
	group, e := _self.repo.Query(&entity.Group{Id: target})
	if e != nil {
		return nil, log.WithError(utils.ERR_GROUP_GET_FAIL)
	}
	if group == nil || refresh {
		resultDTO, err := Post("/api/group/selectOne", map[string]uint64{"id": target})
		if err != nil {
			return nil, log.WithError(err)
		}
		//如果服务器获取失败
		if resultDTO.Data == nil {
			return nil, log.WithError(utils.ERR_GROUP_GET_FAIL)
		}
		var g entity.Group
		e := util.Str2Obj(resultDTO.Data.(string), &g)
		if e != nil {
			return nil, log.WithError(utils.ERR_GROUP_GET_FAIL)
		}
		if g.Id != 0 {
			//保存到数据库
			g.UserId = conf.GetLoginInfo().User.Id
			e := _self.repo.Save(&g)
			if e != nil {
				return nil, log.WithError(utils.ERR_OPERATION_FAIL)
			}
		}
		return &g, nil
	}
	return group, nil
}

func (_self *GroupService) NeedPassword(tp string, target uint64) string {
	if "group" != tp {
		return "-1"
	}
	group, e := _self.SelectOne(target, false)
	if e != nil || group == nil {
		return "-1"
	}
	if group.Type == 2 {
		//需要密码 但是有密码
		if conf.Conf.Pwds[tp+"_"+util.Uint642Str(target)] != "" {
			return "2"
		}
		return "1" //前端 仅在返回1时 打开密码输入框
	}
	return "-1"
}

func (_self *GroupService) Update(id uint64, data string, updateType int) *utils.Error {
	//修改本地的群名称
	group, e := _self.repo.Query(&entity.Group{Id: id})
	if e != nil {
		log.Error(e)
		return log.WithError(utils.ERR_GROUP_GET_FAIL)
	}
	switch updateType {
	case 1:
		_, err := Post("/api/group/edit", map[string]interface{}{"id": id, "updateType": updateType, "name": data})
		if err != nil {
			return log.WithError(err)
		}
		group.Name = data
		//修改聊天名称
		chat, e := NewChatService().repo.Query(&entity.Chat{Type: "group", TargetId: id, UserId: conf.GetLoginInfo().User.Id})
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_GROUP_GET_FAIL)
		}
		chat.Name = data
		e = NewChatService().repo.Save(chat)
		break
	case 2:
		_, err := Post("/api/group/edit", map[string]interface{}{"id": id, "updateType": updateType, "notice": data})
		if err != nil {
			return log.WithError(err)
		}
		break
	case 3:
		_, err := Post("/api/group/edit", map[string]interface{}{"id": id, "updateType": updateType, "headImg": data})
		if err != nil {
			return log.WithError(err)
		}
		break
	}
	e = _self.repo.Save(group)
	if e != nil {
		log.Error(e)
		return log.WithError(utils.ERR_GROUP_GET_FAIL)
	}
	return nil
}

func (_self *GroupService) Quit(id uint64) *utils.Error {
	tx := _self.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		err := _self.DelLocalGroup(id)
		if err != nil {
			return log.WithError(err)
		}
		_, err = Post("/api/group/quit", map[string]interface{}{"id": id})
		if err != nil {
			return log.WithError(err)
		}
		e := tx.Commit().Error
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

func (_self *GroupService) Delete(id uint64) *utils.Error {
	tx := _self.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		err := _self.DelLocalGroup(id)
		if err != nil {
			return log.WithError(err)
		}
		_, err = Post("/api/group/delete", map[string]interface{}{"id": id})
		if err != nil {
			return log.WithError(err)
		}
		e := tx.Commit().Error
		if e != nil {
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		log.Debugf("提交事务33333")
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}

func (_self *GroupService) DelLocalGroup(id uint64) *utils.Error {
	tx := _self.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		log.Debugf("del local group %d", id)
		//删除群
		e := _self.repo.Delete(&entity.Group{
			Id: id,
		})
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		log.Debugf("del local group member %d", id)
		//删除群成员
		e = NewGroupMemberService().repo.Delete(&entity.GroupMember{
			GId: id,
		})
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		log.Debugf("del local chat %d", id)
		//删除聊天
		e = NewChatService().DelLocalChat("group", id)
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}

		e = tx.Commit().Error
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_DEL_FAIL)
		}
		log.Debugf("提交事务成功22222")
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// SelectOneGroupMemberInfo 获取群成员信息
func (_self *GroupService) SelectOneGroupMemberInfo(gId, userId uint64) (map[string]interface{}, *utils.Error) {
	//如果是当前用户 就什么都不返回
	if userId == conf.GetLoginInfo().User.Id {
		return nil, nil
	}
	//获取用户信息
	user, err := NewUserService().SelectOne(userId, false)
	if err != nil {
		log.Error(err)
		return nil, log.WithError(utils.ERR_GET_USER_FAIL)
	}
	//用户信息不存在 就查一次
	if user == nil {
		user, err = NewUserService().SelectOne(userId, true)
		if err != nil {
			log.Error(err)
			return nil, log.WithError(utils.ERR_GET_USER_FAIL)
		}
	}
	//获取群成员信息
	gm, e := NewGroupMemberService().repo.Query(&entity.GroupMember{
		GId:    gId,
		UserId: userId,
		State:  2,
	})
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_GET_USER_FAIL)
	}
	data := map[string]interface{}{
		"userId":  userId,
		"headImg": user.HeadImg,
	}
	//有群昵称
	if gm != nil && gm.Name != "" {
		data["name"] = gm.Name
		return data, nil
	}
	//没有群昵称 获取好友信息
	friend, err := NewFriendService().IsFriend(userId)
	if err != nil {
		log.Error(err)
		return nil, log.WithError(utils.ERR_GET_USER_FAIL)
	}
	if friend != nil && friend.Name != "" {
		data["name"] = friend.Name
		return data, nil
	}
	//没有好友信息
	data["name"] = user.Nickname
	return data, nil
}
