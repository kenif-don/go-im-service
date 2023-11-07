package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type ChatRepo struct {
	*db.Transaction
}

func NewChatRepo() *ChatRepo {
	return &ChatRepo{Transaction: db.NewTransaction()}
}
func (_self *ChatRepo) Query(obj *entity.Chat) (*entity.Chat, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *ChatRepo) QueryAll(obj *entity.Chat) ([]entity.Chat, error) {
	objs := &[]entity.Chat{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *ChatRepo) Save(obj *entity.Chat) error {
	tx := _self.Data.Db.Where(obj.Id).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *ChatRepo) Delete(obj *entity.Chat) error {
	tx := _self.Data.Db.Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
