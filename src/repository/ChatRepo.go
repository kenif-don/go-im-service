package repository

import (
	"errors"
	"go-im-service/src/configs/db"
	"go-im-service/src/entity"

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
func (_self *ChatRepo) QueryAll(userId uint64) ([]entity.Chat, error) {
	objs := &[]entity.Chat{}
	tx := _self.Data.Db.Where("user_id = ?", userId).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *ChatRepo) Save(obj *entity.Chat) error {
	chat, e := _self.Query(&entity.Chat{
		Type:     obj.Type,
		TargetId: obj.TargetId,
		UserId:   obj.UserId,
	})
	if e != nil {
		return e
	}
	if chat == nil {
		tx := _self.Data.Db.Create(obj)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	}
	chat.Name = obj.Name
	chat.HeadImg = obj.HeadImg
	chat.UnReadNo = obj.UnReadNo
	chat.Top = obj.Top
	tx := _self.Data.Db.Model(&entity.Chat{}).Where(chat.Id).Save(chat)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *ChatRepo) Delete(obj *entity.Chat) error {
	tx := _self.Data.Db.Model(&entity.Chat{}).Where("`type` = ? and `target_id` = ? and `user_id` = ?", obj.Type, obj.TargetId, obj.UserId).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
