package service

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"gorm.io/gorm"
)

type IFriendRepo interface {
	Query(obj *entity.Friend) (*entity.Friend, error)
	QueryAll(obj *entity.Friend) (*[]entity.Friend, error)
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
func QueryFriend(id uint64, repo IFriendRepo) (*entity.Friend, error) {
	return repo.Query(&entity.Friend{Id: id})
}
func QueryFriendAll(repo IFriendRepo) (*[]entity.Friend, error) {
	return repo.QueryAll(&entity.Friend{})
}
func (_self *FriendService) UpdateOne(he, me uint64) *utils.Error {
	var req = make(map[string]uint64)
	req["he"] = he
	req["me"] = me
	resultDTO, err := util.Post("/api/friend/selectOne", req)
	if err != nil {
		return log.WithError(err)
	}
	var fa entity.Friend
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
func (_self *FriendService) UpdateName(id uint64, name string) *utils.Error {
	friend, e := QueryFriend(id, _self.repo)
	if e != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}

	e = _self.repo.Save(friend)
	if e != nil {
		return log.WithError(utils.ERR_OPERATION_FAIL)
	}
	return nil
}
