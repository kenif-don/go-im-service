package entity

type Chat struct {
	Id       uint64 `gorm:"unique;<-:create" json:"id"`
	Type     string `json:"type" gorm:"index:idx_2"`     // 聊天类型 friend, group
	TargetId uint64 `json:"targetId" gorm:"index:idx_2"` // 聊天目标 用户ID或群ID
	UserId   uint64 `json:"userId" gorm:"index:idx_2"`   // 当前聊天所有者 用户ID
	Name     string `json:"name"`                        // 聊天名称
	HeadImg  string `json:"headImg"`                     // 聊天头像
	UnReadNo int    `json:"unRead"`                      // 未读消息数量
	Top      int    `json:"top"`                         // 是否置顶 1:置顶 0:不置顶
	LastMsg  string `gorm:"-" json:"lastMsg"`            // 最后一条聊天
	LastTime uint64 `gorm:"-" json:"time"`               // 最后一条聊天时间
}
