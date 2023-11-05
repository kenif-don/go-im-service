package service

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
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
