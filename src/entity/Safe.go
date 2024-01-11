package entity

type Safe struct {
	Id         uint64 `json:"id"`
	Content    string `json:"content"`
	UserId     uint64 `json:"userId"`
	CreateTime uint64 `json:"createTime"`
}
