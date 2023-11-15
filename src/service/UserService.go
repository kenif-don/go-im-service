package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/db"
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

// UpdatePassword 修改密码 修改后需要推送重新登录通知
func (_self *UserService) UpdatePassword(tp int, pwd, oldPwd, newPwd string) *utils.Error {
	resultDTO, err := Post("/api/user/editPwd", map[string]interface{}{"type": tp, "pwd": pwd, "oldPwd": oldPwd, "newPwd": newPwd})
	if err != nil {
		return log.WithError(err)
	}
	user := &entity.User{}
	e := util.Str2Obj(resultDTO.Data.(string), user)
	if e != nil {
		return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
	}
	//查找到数据库中存的私钥--设置到服务器返回的对象中
	sysUser, e := QueryUser(conf.GetLoginInfo().User.Id, _self.repo)
	if e != nil || user == nil {
		return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
	}
	user.PrivateKey = sysUser.PrivateKey
	e = _self.repo.Save(user)
	if e != nil {
		return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
	}
	if tp == 1 {
		// 重新登录
		err = _self.Logout()
		if err != nil {
			return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
		}
	}
	return nil
}
func (_self *UserService) UpdateUser(id uint64) (*entity.User, *utils.Error) {
	resultDTO, err := Post("/api/user/selectOne", map[string]interface{}{"id": id})
	if err != nil {
		return nil, log.WithError(utils.ERR_USER_UPDATE_FAIL)
	}
	var user = &entity.User{}
	e := util.Str2Obj(resultDTO.Data.(string), user)
	if e != nil {
		return nil, log.WithError(utils.ERR_USER_UPDATE_FAIL)
	}
	e = _self.repo.Save(user)
	if e != nil {
		return nil, log.WithError(utils.ERR_USER_UPDATE_FAIL)
	}
	return user, nil
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
	if len(intro) < 1 || len(intro) > 10 {
		return log.WithError(utils.ERR_INTRO_VALIDATE_FAIL)
	}
	user, err := QueryUser(id, _self.repo)
	if err != nil || user == nil {
		return log.WithError(utils.ERR_INTRO_UPDATE_FAIL)
	}
	user.Intro = intro
	return _self.Update(user)
}
func (_self *UserService) UpdateNickname(id uint64, nickname string) *utils.Error {
	if len(nickname) < 1 || len(nickname) > 10 {
		return log.WithError(utils.ERR_NICKNAME_VALIDATE_FAIL)
	}
	user, err := QueryUser(id, _self.repo)
	if err != nil || user == nil {
		return log.WithError(utils.ERR_NICKNAME_UPDATE_FAIL)
	}
	user.Nickname = nickname
	return _self.Update(user)
}
func (_self *UserService) Update(obj *entity.User) *utils.Error {
	tx := _self.repo.BeginTx()
	if e := tx.Error; e != nil {
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
	if len(username) < 6 || len(username) > 20 {
		return log.WithError(utils.ERR_USER_USERNAME_LENGTH)
	}
	if len(password) < 6 || len(password) > 20 {
		return log.WithError(utils.ERR_USER_PASSWORD_LENGTH)
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
	if len(username) < 6 || len(username) > 20 {
		return log.WithError(utils.ERR_USER_USERNAME_LENGTH)
	}
	if len(password) < 6 || len(password) > 20 {
		return log.WithError(utils.ERR_USER_PASSWORD_LENGTH)
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
func (_self *UserService) LoginPwd2(pwd2 string) *utils.Error {
	if len(pwd2) < 6 || len(pwd2) > 20 {
		return log.WithError(utils.ERR_USER_PASSWORD_LENGTH)
	}
	//远程服务器确认密码是否正确
	resultDTO, err := Post("/api/user/loginPwd2", map[string]string{"password2": pwd2})
	if err != nil {
		return log.WithError(err)
	}
	//如果自毁 删聊天记录、删聊天
	if resultDTO.Data != nil && resultDTO.Data.(string) == "burst" {
		tx := db.NewTransaction().BeginTx()
		if e := tx.Error; e != nil {
			return log.WithError(utils.ERR_QUERY_FAIL)
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err = func() *utils.Error {
			err = NewChatService().DelAllChat()
			if err != nil {
				return log.WithError(err)
			}
			err = NewMessageService().DelAllMessage()
			if err != nil {
				return log.WithError(err)
			}
			return nil
		}()
		if err != nil {
			tx.Rollback()
		}
		return log.WithError(err)
	}
	//没有错误 标记已经输入2级密码
	conf.UpdateInputPwd2(2)
	return nil
}
func (_self *UserService) Logout() *utils.Error {
	conf.ClearLoginInfo()
	mgr := handler.GetClientHandler().GetMessageManager()
	if mgr == nil {
		return log.WithError(utils.ERR_NET_FAIL)
	}
	mgr.SendLogout()
	//通知前往登录页面
	if Listener != nil {
		Listener.OnLogin()
	}
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
