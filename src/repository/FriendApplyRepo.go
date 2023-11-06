package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type FriendApplyRepo struct {
	*db.Transaction
}

func NewFriendApplyRepo() *FriendApplyRepo {
	return &FriendApplyRepo{Transaction: db.NewTransaction()}
}
func (_self *FriendApplyRepo) Query(obj *entity.FriendApply) (*entity.FriendApply, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *FriendApplyRepo) QueryAll(obj *entity.FriendApply) ([]entity.FriendApply, error) {
	objs := &[]entity.FriendApply{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *FriendApplyRepo) Save(obj *entity.FriendApply) error {
	tx := _self.Data.Db.Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *FriendApplyRepo) Delete(obj *entity.FriendApply) error {
	tx := _self.Data.Db.Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
