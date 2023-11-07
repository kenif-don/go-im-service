package entity

type Chat struct {
	Id        uint64    `gorm:"unique;<-:create" json:"id"`
	Type      string    `json:"type"`               // 聊天类型 friend, group
	TargetId  uint64    `json:"targetId"`           // 聊天目标 好友ID或群ID
	UserId    uint64    `json:"userId"`             // 当前聊天所有者 用户ID
	Name      string    `json:"name"`               // 聊天名称
	HeadImg   string    `json:"headImg"`            // 聊天头像
	UnReadNo  int       `json:"unRead"`             // 未读消息数量
	LastMsg   string    `gorm:"-" json:"msg"`       // 最后一条聊天
	LastTime  uint64    `gorm:"-" json:"time"`      // 最后一条聊天时间
	Msgs      []Message `gorm:"-" json:"msgs"`      // 分页消息
	Page      int       `gorm:"-" json:"page"`      // 当前页码
	TotalPage int       `gorm:"-" json:"totalPage"` // 总页码
}
