package entity

type RechargeOrder struct {
	Id         uint64 `gorm:"unique;<-:create" json:"id"`
	UserId     uint64 `gorm:"unique" json:"userId"` //订单
	Value      string `json:"value"`                //重置金额
	Status     int    `json:"status"`
	Type       int    `json:"type"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
	Wallet     string `json:"wallet"`
}
