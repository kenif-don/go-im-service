package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type GroupMemberRepo struct {
	*db.Transaction
}

func NewGroupMemberRepo() *GroupMemberRepo {
	return &GroupMemberRepo{Transaction: db.NewTransaction()}
}
func (_self *GroupMemberRepo) Query(obj *entity.GroupMember) (*entity.GroupMember, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *GroupMemberRepo) QueryAll(obj *entity.GroupMember) ([]entity.GroupMember, error) {
	objs := &[]entity.GroupMember{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *GroupMemberRepo) Save(obj *entity.GroupMember) error {
	tx := _self.Data.Db.Model(&entity.GroupMember{}).Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *GroupMemberRepo) Delete(obj *entity.GroupMember) error {
	tx := _self.Data.Db.Model(&entity.GroupMember{}).Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
