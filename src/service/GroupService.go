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
	resultDTO, err := Post("/api/group/create", map[string]interface{}{"ids": ids, "tp": tp, "password": password})
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
func (_self *GroupService) QueryAll() ([]entity.Group, error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return []entity.Group{}, nil
	}
	return _self.repo.QueryAll(&entity.Group{Owner: conf.GetLoginInfo().User.Id})
}
