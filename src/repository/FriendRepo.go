package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type FriendRepo struct {
	*db.Transaction
}

func NewFriendRepo() *FriendRepo {
	return &FriendRepo{Transaction: db.NewTransaction()}
}
func (_self *FriendRepo) Query(obj *entity.Friend) (*entity.Friend, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *FriendRepo) QueryAll(obj *entity.Friend) (*[]entity.Friend, error) {
	objs := &[]entity.Friend{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return objs, nil
}
func (_self *FriendRepo) Save(obj *entity.Friend) error {
	tx := _self.Data.Db.Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *FriendRepo) Delete(obj *entity.Friend) error {
	tx := _self.Data.Db.Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
