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

type IFriendApplyRepo interface {
	Query(obj *entity.FriendApply) (*entity.FriendApply, error)
	QueryAll(obj *entity.FriendApply) (*[]entity.FriendApply, error)
	Save(obj *entity.FriendApply) error
	Delete(obj *entity.FriendApply) error
	BeginTx() *gorm.DB
}
type FriendApplyService struct {
	repo IFriendApplyRepo
}

func NewFriendApplyService() *FriendApplyService {
	return &FriendApplyService{
		repo: repository.NewFriendApplyRepo(),
	}
}
func QueryFriendApply(fa *entity.FriendApply, repo IFriendApplyRepo) (*entity.FriendApply, error) {
	return repo.Query(fa)
}

// QueryFriendApplyAll 查询登录者的所有好友请求
func QueryFriendApplyAll(repo IFriendApplyRepo) (*[]entity.FriendApply, error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return &[]entity.FriendApply{}, nil
	}
	return repo.QueryAll(&entity.FriendApply{To: conf.GetLoginInfo().User.Id})
}

// UpdateOne 查询单个 然后同步到数据库
func (_self *FriendApplyService) UpdateOne(obj *entity.FriendApply) *utils.Error {
	//没用查到 就从后台查一次
	resultDTO, err := util.Post("/api/friend-apply/selectOne", obj)
	if err != nil {
		return log.WithError(err)
	}
	var fa entity.FriendApply
	_ = util.Str2Obj(resultDTO.Data.(string), &fa)
	if fa.Id != 0 {
		//保存到数据库
		e := _self.repo.Save(&fa)
		if e != nil {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
	}
	return nil
}

// SelectAll 查询登录者的所有好友请求 别人请求自己的
func (_self *FriendApplyService) SelectAll() (*[]entity.FriendApply, *utils.Error) {
	//先从数据库查
	sysFriendApplys, e := QueryFriendApplyAll(_self.repo)
	if e != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	if len(*sysFriendApplys) != 0 {
		//封装用户信息
		for _, v := range *sysFriendApplys {
			userService := NewUserService()
			user, e := QueryUser(v.To, userService.repo)
			if e != nil || user == nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			v.FromUser = user
		}
		return sysFriendApplys, nil
	}
	//没用查到 就从后台查一次
	resultDTO, err := util.Post("/api/friend-apply/selectAll", nil)
	if err != nil {
		return nil, log.WithError(err)
	}
	var fas []entity.FriendApply
	e = util.Str2Obj(resultDTO.Data.(string), &fas)
	if e != nil || fas == nil {
		return &[]entity.FriendApply{}, nil
	}
	tx := _self.repo.BeginTx()
	if err := tx.Error; err != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	r, err := func() (*[]entity.FriendApply, *utils.Error) {
		// 保存到数据库
		for _, v := range fas {
			e := _self.repo.Save(&v)
			if e != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			//先查询 是否存在 存在就不添加了
			//保存对应的用户信息
			userService := NewUserService()
			sysUser, e := QueryUser(v.To, userService.repo)
			if e != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
			if sysUser != nil {
				continue
			}
			e = userService.Save(v.FromUser)
			if e != nil {
				return nil, log.WithError(utils.ERR_QUERY_FAIL)
			}
		}
		e = tx.Commit().Error
		if e != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		return &fas, nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return r, err
}

// Update 拒绝或同意
func (_self *FriendApplyService) Update(id uint64, state int) *utils.Error {
	if id == 0 {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	tx := _self.repo.BeginTx()
	if err := tx.Error; err != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		sysFriendApply, e := QueryFriendApply(&entity.FriendApply{Id: id}, _self.repo)
		if e != nil || sysFriendApply == nil {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		if sysFriendApply.State != 1 {
			return log.WithError(utils.ERR_OPERATION_FAIL)
		}
		sysFriendApply.State = state
		_, err := util.Post("/api/friend-apply/edit", sysFriendApply)
		if err != nil {
			return log.WithError(err)
		}
		// 请求存到数据库
		e = _self.repo.Save(sysFriendApply)
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

// Add 发起添加好友请求 自己发的请求不会添加到数据库
func (_self *FriendApplyService) Add(to uint64, remark string) *utils.Error {
	if to == 0 {
		return log.WithError(utils.ERR_ADD_FRIEND_FAIL)
	}
	//发起添加请求
	obj := &entity.FriendApply{
		To:     to,
		Remark: remark,
	}
	_, err := util.Post("/api/friend-apply/add", obj)
	if err != nil {
		return log.WithError(err)
	}
	return nil
}
