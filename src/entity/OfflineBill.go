package entity

type OfflineBill struct {
	Id      int64  `json:"id"`
	UserId  int64  `json:"user_id"`
	Content string `json:"content"`
}
