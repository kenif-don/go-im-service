package db

import (
	"gorm.io/gorm"
)

var Tx *Transaction

type Transaction struct {
	Data *DB
}

func (_self *Transaction) BeginTx() *gorm.DB {
	return _self.Data.Db.Begin().WithContext(Ctx)
}

func NewTransaction() *Transaction {
	if Tx != nil {
		return Tx
	}
	data := NewDB().Db.WithContext(Ctx)
	return &Transaction{
		Data: &DB{Db: data},
	}
}
