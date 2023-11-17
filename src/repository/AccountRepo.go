package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type AccountRepo struct {
	*db.Transaction
}

func NewAccountRepo() *AccountRepo {
	return &AccountRepo{Transaction: db.NewTransaction()}
}
func (_self *AccountRepo) Query(obj *entity.Account) (*entity.Account, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *AccountRepo) QueryAll(obj *entity.Account) ([]entity.Account, error) {
	objs := &[]entity.Account{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *AccountRepo) Save(obj *entity.Account) error {
	tx := _self.Data.Db.Model(&entity.Account{}).Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *AccountRepo) Delete(obj *entity.Account) error {
	tx := _self.Data.Db.Model(&entity.Account{}).Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
