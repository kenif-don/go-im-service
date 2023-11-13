package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"gorm.io/gorm"
	"im-sdk/handler"
)

type IUserRepo interface {
	Query(obj *entity.User) (*entity.User, error)
	QueryAll(obj *entity.User) (*[]entity.User, error)
	Save(obj *entity.User) error
	Delete(obj *entity.User) error
	BeginTx() *gorm.DB
}
type UserService struct {
	repo IUserRepo
}

func NewUserServiceNoDB() *UserService {
	return &UserService{}
}
func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepo(),
	}
}
func QueryUser(id uint64, repo IUserRepo) (*entity.User, error) {
	return repo.Query(&entity.User{Id: id})
}
func (_self *UserService) Save(obj *entity.User) error {
	return _self.repo.Save(obj)
}
func (_self *UserService) Search(keyword string) (string, *utils.Error) {
	if keyword == "" {
		return "", nil
	}
	var data = make(map[string]string)
	data["keyword"] = keyword
	resultDTO, err := Post("/api/user/search", data)
	if err != nil {
		return "", err
	}
	return resultDTO.Data.(string), nil
}

//	func (_self *UserService) UpdatePassword(id uint64, password string) *utils.Error {
//		user, err := QueryUser(id, _self.repo)
//		if err != nil || user == nil {
//			return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
//		}
//		user.Password = password
//		return _self.Update(user)
//	}
func (_self *UserService) UpdateTargetPublicKey(id uint64) *utils.Error {
	user, err := QueryUser(id, _self.repo)
	if err != nil || user == nil {
		return log.WithError(utils.ERR_HEADIMG_UPDATE_FAIL)
	}
	resultDTO, e2 := Post("/api/user/GetTargetPublicKey", map[string]interface{}{"id": id})
	if e2 != nil {
		return log.WithError(utils.ERR_SEND_FAIL)
	}
	user.PublicKey = resultDTO.Data.(string)
	e := _self.repo.Save(user)
	if e != nil {
		return log.WithError(utils.ERR_SEND_FAIL)
	}
	return nil
}
func (_self *UserService) UpdateHeadImg(id uint64, headImg string) *utils.Error {
	user, err := QueryUser(id, _self.repo)
	if err != nil || user == nil {
		return log.WithError(utils.ERR_HEADIMG_UPDATE_FAIL)
	}
	user.HeadImg = headImg
	return _self.Update(user)
}
func (_self *UserService) UpdateEmail(id uint64, email string) *utils.Error {
	user, err := QueryUser(id, _self.repo)
	if err != nil || user == nil {
		return log.WithError(utils.ERR_EMAIL_UPDATE_FAIL)
	}
	user.Email = email
	return _self.Update(user)
}
func (_self *UserService) UpdateIntro(id uint64, intro string) *utils.Error {
	user, err := QueryUser(id, _self.repo)
	if err != nil || user == nil {
		return log.WithError(utils.ERR_INTRO_UPDATE_FAIL)
	}
	user.Intro = intro
	return _self.Update(user)
}
func (_self *UserService) UpdateNickname(id uint64, nickname string) *utils.Error {
	user, err := QueryUser(id, _self.repo)
	if err != nil || user == nil {
		return log.WithError(utils.ERR_NICKNAME_UPDATE_FAIL)
	}
	user.Nickname = nickname
	return _self.Update(user)
}
func (_self *UserService) Update(obj *entity.User) *utils.Error {
	tx := _self.repo.BeginTx()
	if err := tx.Error; err != nil {
		return log.WithError(utils.ERR_USER_UPDATE_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		//发起请求修改后台用户信息
		_, err := Post("/api/user/edit", obj)
		if err != nil {
			return log.WithError(err)
		}
		//修改数据库
		e := _self.repo.Save(obj)
		if e != nil {
			return log.WithError(utils.ERR_USER_UPDATE_FAIL)
		}
		//修改登录者
		conf.PutLoginInfo(*obj)
		e = tx.Commit().Error
		if e != nil {
			return log.WithError(utils.ERR_USER_UPDATE_FAIL)
		}
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// UpdateLoginUserKeys 修改公私钥
func (_self *UserService) UpdateLoginUserKeys(keys util.EncryptKeys) *utils.Error {
	obj, e := QueryUser(conf.LoginInfo.User.Id, _self.repo)
	if e != nil {
		return log.WithError(utils.ERR_SECRET_UPDATE_FAIL)
	}
	if obj == nil {
		return log.WithError(utils.ERR_SECRET_UPDATE_FAIL)
	}
	obj.PublicKey = keys.PublicKey
	obj.PrivateKey = keys.PrivateKey
	tx := _self.repo.BeginTx()
	if e := tx.Error; e != nil {
		return log.WithError(utils.ERR_SECRET_UPDATE_FAIL)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := func() *utils.Error {
		//发起请求修改后台用户信息
		_, err := Post("/api/user/resetPublicKey", &entity.User{PublicKey: keys.PublicKey})
		if err != nil {
			return log.WithError(err)
		}
		//修改本地信息
		e = _self.repo.Save(obj)
		if e != nil {
			return log.WithError(utils.ERR_SECRET_UPDATE_FAIL)
		}
		//修改登录者
		conf.PutLoginInfo(*obj)
		e = tx.Commit().Error
		if e != nil {
			return log.WithError(utils.ERR_SECRET_UPDATE_FAIL)
		}
		return nil
	}()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// Register 用户注册
func (_self *UserService) Register(username, password string) *utils.Error {
	if username == "" {
		return log.WithError(utils.ERR_USER_REGISTER_USERNAME_NULL)
	}
	if password == "" {
		return log.WithError(utils.ERR_USER_REGISTER_PASSWORD_NULL)
	}
	if len(username) < 6 || len(username) > 20 {
		return log.WithError(utils.ERR_USER_REGISTER_USERNAME_LENGTH)
	}
	if len(password) < 6 || len(password) > 20 {
		return log.WithError(utils.ERR_USER_REGISTER_PASSWORD_LENGTH)
	}
	params := &entity.RegisterUser{
		Username: username,
		Password: password,
	}
	_, err := Post("/api/user/register", params)
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

func (_self *UserService) Login(username, password string) *utils.Error {
	if username == "" {
		return log.WithError(utils.ERR_USER_REGISTER_USERNAME_NULL)
	}
	if password == "" {
		return log.WithError(utils.ERR_USER_REGISTER_PASSWORD_NULL)
	}
	if len(username) < 6 || len(username) > 20 {
		return log.WithError(utils.ERR_USER_REGISTER_USERNAME_LENGTH)
	}
	if len(password) < 6 || len(password) > 20 {
		return log.WithError(utils.ERR_USER_REGISTER_PASSWORD_LENGTH)
	}

	params := &entity.RegisterUser{
		Username: username,
		Password: password,
	}
	resultDTO, err := Post("/api/user/login", params)
	if err != nil || resultDTO == nil || resultDTO.Data == nil {
		return log.WithError(err)
	}
	//缓存登录token
	conf.PutToken(resultDTO.Data.(string))
	// 不知道是否需要输入
	conf.UpdateInputPwd2(1)
	//获取用户信息
	return _self.LoginInfo()
}
func (_self *UserService) Logout() *utils.Error {
	conf.ClearLoginInfo()
	mgr := handler.GetClientHandler().GetMessageManager()
	if mgr == nil {
		return log.WithError(utils.ERR_NET_FAIL)
	}
	mgr.SendLogout()
	return nil
}

// LoginInfo 获取用户信息
func (_self *UserService) LoginInfo() *utils.Error {
	resultDTO, err := Post("/api/user/info", nil)
	if err != nil || resultDTO == nil || resultDTO.Data == nil {
		return log.WithError(err)
	}
	//缓存用户信息
	var user entity.User
	e := util.Obj2Obj(resultDTO.Data, &user)
	if e != nil {
		return log.WithError(utils.ERR_GET_USER_INFO)
	}
	//数据库不存在 就添加 这里不做修改
	sysUser, e := QueryUser(user.Id, _self.repo)
	if e != nil {
		return log.WithError(utils.ERR_GET_USER_INFO)
	}
	//数据存在--需要把数据库中的私钥封装到登录者中
	if sysUser != nil {
		user.PrivateKey = sysUser.PrivateKey
		//存到文件--如果没有 会重新生成公私钥
		conf.PutLoginInfo(user)
		return nil
	}
	e = _self.Save(&user)
	if e != nil {
		return log.WithError(utils.ERR_GET_USER_INFO)
	}
	//存到文件--如果没有 会重新生成公私钥
	conf.PutLoginInfo(user)
	return nil
}
