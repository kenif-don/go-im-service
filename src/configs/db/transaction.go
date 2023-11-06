package db

import (
	"gorm.io/gorm"
)

type Transaction struct {
	Data *DB
}

func (_self *Transaction) BeginTx() *gorm.DB {
	//data := _self.Data.Db.Begin().WithContext(Ctx)
	//_self.Data = &DB{Db: data}
	//return data
	return _self.Data.Db.Begin()
}

func NewTransaction() *Transaction {
	data := NewDB().Db.WithContext(Ctx)
	return &Transaction{
		Data: &DB{Db: data},
	}
}
