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

type IGroupMemberRepo interface {
	Query(obj *entity.GroupMember) (*entity.GroupMember, error)
	QueryAll(obj *entity.GroupMember) ([]entity.GroupMember, error)
	Save(obj *entity.GroupMember) error
	Delete(obj *entity.GroupMember) error
	BeginTx() *gorm.DB
}
type GroupMemberService struct {
	repo IGroupMemberRepo
}

func NewGroupMemberService() *GroupMemberService {
	return &GroupMemberService{
		repo: repository.NewGroupMemberRepo(),
	}
}

// SelectMembers 从服务器获取群成员
func (_self *GroupMemberService) SelectMembers(gId uint64, refresh bool) ([]entity.GroupMember, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	ms, e := _self.QueryAll(&entity.GroupMember{GId: gId})
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	if ms != nil && len(ms) > 0 && !refresh {
		return ms, nil
	}
	//需要刷新 先删除一次 重新获取
	if refresh {
		e = _self.repo.Delete(&entity.GroupMember{GId: gId})
		if e != nil {
			log.Error(e)
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
	}
	resultDTO, err := Post("/api/group/selectMembers", map[string]interface{}{"id": gId})
	if err != nil {
		return nil, log.WithError(err)
	}
	var members []entity.GroupMember
	e = util.Str2Obj(resultDTO.Data.(string), &members)
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//遍历members
	for i := 0; i < len(members); i++ {
		if members[i].User == nil || members[i].User.Id == 0 {
			continue
		}
		//从服务器拉去用户信息
		_, err := NewUserService().SelectOne(members[i].UserId, false)
		if err != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		// 保存群成员
		e := _self.repo.Save(&members[i])
		if e != nil {
			log.Error(e)
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
	}
	return members, nil
}
func (_self *GroupMemberService) QueryAll(gm *entity.GroupMember) ([]entity.GroupMember, error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return []entity.GroupMember{}, nil
	}
	gm.State = 2
	return _self.repo.QueryAll(gm)
}

func (_self *GroupMemberService) GetGroupMemberByMe(id uint64) (*entity.GroupMember, *utils.Error) {
	gm, e := _self.repo.Query(&entity.GroupMember{GId: id, UserId: conf.GetLoginInfo().User.Id})
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	return gm, nil
}

func (_self *GroupMemberService) UpdateGroupMemberName(id uint64, name string) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return log.WithError(utils.ERR_NICKNAME_UPDATE_FAIL)
	}
	_, err := Post("/api/group/updateGroupMemberName", map[string]interface{}{"gId": id, "name": name})
	if err != nil {
		return log.WithError(utils.ERR_NICKNAME_UPDATE_FAIL)
	}
	gm, e := _self.repo.Query(&entity.GroupMember{GId: id, UserId: conf.GetLoginInfo().User.Id})
	if e != nil {
		return log.WithError(utils.ERR_NICKNAME_UPDATE_FAIL)
	}
	gm.Name = name
	e = _self.repo.Save(gm)
	if e != nil {
		return log.WithError(utils.ERR_NICKNAME_UPDATE_FAIL)
	}
	return nil
}
