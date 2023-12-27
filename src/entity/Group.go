package entity

type Group struct {
	Id       uint64 `gorm:"unique;<-:create" json:"id"`
	Name     string `json:"name"`
	Notice   string `json:"notice"`
	HeadImg  string `json:"headImg"`
	Owner    uint64 `json:"owner"`
	State    int    `json:"state"`
	Type     int    `json:"type"`
	Password string `json:"password"`

	HeUser *User `gorm:"-" json:"heUser"`
}
type GroupMember struct {
	Id     uint64 `gorm:"unique;<-:create" json:"id"`
	GId    uint64 `json:"gId"`
	UserId uint64 `json:"userId"`
	Type   int    `json:"type"`
	Name   string `json:"name"`
	State  int    `json:"state"`
	Remark string `json:"remark"`

	User *User `gorm:"-" json:"User"`
}
