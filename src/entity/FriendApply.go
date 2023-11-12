package entity

type FriendApply struct {
	Id         uint64 `gorm:"unique;<-:create" json:"id"`
	From       uint64 `json:"from"`
	To         uint64 `json:"to"`
	Remark     string `json:"remark"`
	State      int    `json:"state"`
	CreateTime int64  `json:"createTime"`

	FromUser *User `gorm:"-" json:"fromUser"`
}
