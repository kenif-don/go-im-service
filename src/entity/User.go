package entity

type User struct {
	Id            uint64 `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Password2     string `json:"password2"`
	BurstPassword string `json:"burstPassword"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	PublicKey     string `json:"publicKey"`
	PrivateKey    string `json:"privateKey"`
	Qrcode        string `json:"qrcode"`
	Intro         string `json:"intro"`
	HeadImg       string `json:"headImg"`
	VersionCode   string `json:"versionCode"`
	NoticeType    int    `json:"noticeType"`
}
type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
