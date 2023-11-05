package entity

type Friend struct {
	Id    uint64 `gorm:"unique;<-:create" json:"id"`
	Me    uint64 `json:"me"`
	He    uint64 `json:"he"`
	Name  string `json:"name"`
	State int    `json:"state"`

	HeUser *User `gorm:"-"`
}
