package db

import (
	"context"
	"gorm.io/gorm"
)

type Transaction struct {
	ctx  context.Context
	Data *DB
}

func (_self *Transaction) BeginTx() *gorm.DB {
	data := _self.Data.Db.Begin().WithContext(_self.ctx)
	_self.Data = &DB{Db: data}

	return data

}

func (_self *Transaction) SetTx(data *gorm.DB) {
	_self.Data = &DB{Db: data}
}

func NewTransaction(ctx context.Context) *Transaction {
	data := NewDB().Db.WithContext(ctx)
	return &Transaction{
		ctx:  ctx,
		Data: &DB{Db: data},
	}
}
