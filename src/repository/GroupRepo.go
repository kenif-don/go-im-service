package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type GroupRepo struct {
	*db.Transaction
}

func NewGroupRepo() *GroupRepo {
	return &GroupRepo{Transaction: db.NewTransaction()}
}
func (_self *GroupRepo) Query(obj *entity.Group) (*entity.Group, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *GroupRepo) QueryAll(obj *entity.Group) ([]entity.Group, error) {
	objs := &[]entity.Group{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *GroupRepo) Save(obj *entity.Group) error {
	group, e := _self.Query(&entity.Group{
		Id: obj.Id,
	})
	if e != nil {
		return e
	}
	//没有就保存
	if group == nil {
		tx := _self.Data.Db.Create(obj)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	}
	//有就修改
	tx := _self.Data.Db.Model(&entity.Group{}).Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *GroupRepo) Delete(obj *entity.Group) error {
	tx := _self.Data.Db.Model(&entity.Group{}).Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
