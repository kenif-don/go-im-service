package repository

import (
	"IM-Service/src/configs/db"
	"IM-Service/src/entity"
	"errors"
	"gorm.io/gorm"
)

type MessageRepo struct {
	*db.Transaction
}

func NewMessageRepo() *MessageRepo {
	return &MessageRepo{Transaction: db.NewTransaction()}
}
func (_self *MessageRepo) Query(obj *entity.Message) (*entity.Message, error) {
	tx := _self.Data.Db.Where(obj).First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *MessageRepo) QueryAll(obj *entity.Message) ([]entity.Message, error) {
	objs := &[]entity.Message{}
	tx := _self.Data.Db.Where(obj).Find(objs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *objs, nil
}
func (_self *MessageRepo) Save(obj *entity.Message) error {
	message, e := _self.Query(&entity.Message{No: obj.No})
	if e != nil {
		return e
	}
	if message == nil {
		tx := _self.Data.Db.Create(obj)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	}

	tx := _self.Data.Db.Where("no = ?", obj.No).Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *MessageRepo) Delete(obj *entity.Message) error {
	tx := _self.Data.Db.Model(&entity.Message{}).Where(obj).Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (_self *MessageRepo) QueryLast(obj *entity.Message) (*entity.Message, error) {
	tx := _self.Data.Db.
		Where("type=? and target_id=? and user_id=?", obj.Type, obj.TargetId, obj.UserId).
		Or("type=? and target_id=? and user_id=?", obj.Type, obj.UserId, obj.UserId).
		Order("time desc").First(obj)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return obj, nil
}
func (_self *MessageRepo) Paging(obj *entity.Message) ([]entity.Message, error) {
	objs := &[]entity.Message{}
	var tx *gorm.DB
	if obj.Time > 0 {
		tx = _self.Data.Db.
			Where("type=? and target_id=? and user_id=? and time < ?", obj.Type, obj.TargetId, obj.UserId, obj.Time).
			Or("type=? and target_id=? and user_id=? and time < ?", obj.Type, obj.UserId, obj.UserId, obj.Time).
			Order("time desc").Limit(15).Find(objs)
	} else {
		tx = _self.Data.Db.
			Where("type=? and target_id=? and user_id=?", obj.Type, obj.TargetId, obj.UserId).
			Or("type=? and target_id=? and user_id=?", obj.Type, obj.UserId, obj.UserId).
			Order("time desc").Limit(15).Find(objs)
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return []entity.Message{}, nil
	}
	if tx == nil || tx.Error != nil {
		return nil, tx.Error
	}
	Reverse(objs)
	return *objs, nil
}

// Reverse 数组倒序
func Reverse(arr *[]entity.Message) {
	var temp entity.Message
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}
func (_self *MessageRepo) Count(obj *entity.Message) (int64, error) {
	var count int64
	tx := _self.Data.Db.Where(obj).Count(&count)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
