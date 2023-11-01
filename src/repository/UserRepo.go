package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserRepo struct {
	*db.Transaction
}

func NewUserRepo(ctx context.Context) *UserRepo {
	return &UserRepo{Transaction: db.NewTransaction(ctx)}
}
func (_self *UserRepo) Query(obj *entity.User) (*entity.User, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
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
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return objs, nil
}
func (_self *UserRepo) Create(obj *entity.User) error {
	tx := _self.Data.Db.Create(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *UserRepo) Update(obj *entity.User) error {
	tx := _self.Data.Db.Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *UserRepo) Delete(obj *entity.User) error {
	tx := _self.Data.Db.Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
