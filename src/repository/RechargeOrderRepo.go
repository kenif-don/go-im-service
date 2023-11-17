package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type RechargeOrderRepo struct {
	*db.Transaction
}

func NewRechargeOrderRepo() *RechargeOrderRepo {
	return &RechargeOrderRepo{Transaction: db.NewTransaction()}
}
func (_self *RechargeOrderRepo) Query(obj *entity.RechargeOrder) (*entity.RechargeOrder, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *RechargeOrderRepo) QueryAll(obj *entity.RechargeOrder) ([]entity.RechargeOrder, error) {
	objs := &[]entity.RechargeOrder{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *RechargeOrderRepo) Save(obj *entity.RechargeOrder) error {
	tx := _self.Data.Db.Model(&entity.RechargeOrder{}).Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *RechargeOrderRepo) Delete(obj *entity.RechargeOrder) error {
	tx := _self.Data.Db.Model(&entity.RechargeOrder{}).Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
