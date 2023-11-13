package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type UserRepo struct {
	*db.Transaction
}

func NewUserRepo() *UserRepo {
	return &UserRepo{Transaction: db.NewTransaction()}
}
func (_self *UserRepo) Query(obj *entity.User) (*entity.User, error) {
	tx := _self.Data.Db.Model(&obj).Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *UserRepo) QueryAll(obj *entity.User) (*[]entity.User, error) {
	objs := &[]entity.User{}
	tx := _self.Data.Db.Model(&obj).Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return objs, nil
}
func (_self *UserRepo) Save(obj *entity.User) error {
	tx := _self.Data.Db.Model(&obj).Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *UserRepo) Delete(obj *entity.User) error {
	tx := _self.Data.Db.Model(&obj).Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
