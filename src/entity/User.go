package entity

type User struct {
	Id            uint64 `gorm:"unique;<-:create" json:"id"`
	Username      string `gorm:"<-:create" json:"username"`
	Password      string `gorm:"column:password" json:"password"`
	Password2     string `gorm:"column:password2" json:"password2"`
	BurstPassword string `gorm:"column:burstPassword" json:"burstPassword"`
	SafePassword  string `gorm:"column:safePassword" json:"safePassword"`
	Nickname      string `gorm:"column:nickname" json:"nickname"`
	Email         string `gorm:"column:email" json:"email"`
	PublicKey     string `gorm:"column:publicKey" json:"publicKey"`
	PrivateKey    string `gorm:"column:privateKey" json:"privateKey"`
	Qrcode        string `gorm:"column:qrcode" json:"qrcode"`
	Intro         string `gorm:"column:intro" json:"intro"`
	HeadImg       string `gorm:"column:headImg" json:"headImg"`
	VersionCode   string `gorm:"column:versionCode" json:"versionCode"`
	NoticeType    int    `gorm:"<-:create;column:noticeType" json:"noticeType"`
}
type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
