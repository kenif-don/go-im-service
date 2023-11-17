package entity

type Account struct {
	Id     uint64 `gorm:"unique;<-:create" json:"id"`
	UserId uint64 `gorm:"unique" json:"userId"` //账户所属者
	Amount string `json:"amount"`               //账户余额
}
