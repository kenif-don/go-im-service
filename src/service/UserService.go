package service

import (
	"go-im-service/src/configs/conf"
	"go-im-service/src/configs/db"
	"go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"go-im-service/src/entity"
	"go-im-service/src/repository"
	"go-im-service/src/util"
	"go-nio-client-sdk/handler"

	"gorm.io/gorm"
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
func (_self *UserService) Save(obj *entity.User) error {
	return _self.repo.Save(obj)
}

// SelectOne 先从本地获取 获取失败或需要刷新 就从服务器获取
func (_self *UserService) SelectOne(id uint64, refresh bool) (*entity.User, *utils.Error) {
	//忽略数据库错误，如果出现错误 那么user应该为nil 直接从服务器获取即可
	user, e := _self.repo.Query(&entity.User{Id: id})
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_GET_USER_FAIL)
	}
	if user == nil || refresh {
		resultDTO, err := Post("/api/user/selectOne", map[string]interface{}{"id": id})
		if err != nil {
			return nil, log.WithError(utils.ERR_GET_USER_FAIL)
		}
		if resultDTO.Data == nil {
			return nil, log.WithError(utils.ERR_GET_USER_FAIL)
		}
		user = &entity.User{}
		e := util.Str2Obj(resultDTO.Data.(string), user)
		if e != nil {
			return nil, log.WithError(utils.ERR_GET_USER_FAIL)
		}
		e = _self.repo.Save(user)
		if e != nil {
			return nil, log.WithError(utils.ERR_GET_USER_FAIL)
		}
	}
	return user, nil
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
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string), nil
}

// UpdatePassword 修改密码 修改后需要推送重新登录通知
func (_self *UserService) UpdatePassword(tp int, pwd, oldPwd, newPwd string) *utils.Error {
	resultDTO, err := Post("/api/user/editPwd", map[string]interface{}{"type": tp, "pwd": pwd, "oldPwd": oldPwd, "newPwd": newPwd})
	if err != nil {
		return log.WithError(err)
	}
	if resultDTO.Data == nil {
		return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
	}
	user := &entity.User{}
	e := util.Str2Obj(resultDTO.Data.(string), user)
	if e != nil {
		return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
	}
	//查找到数据库中存的私钥--设置到服务器返回的对象中
	sysUser, err := _self.SelectOne(conf.GetLoginInfo().User.Id, false)
	if err != nil {
		return log.WithError(utils.ERR_PASSWORD_UPDATE_FAIL)
	}
	user.PrivateKey = sysUser.PrivateKey
	conf.PutLoginInfo(*user)
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

func (_self *UserService) UpdateHeadImg(id uint64, headImg string) *utils.Error {
	user, err := _self.SelectOne(id, false)
	if err != nil {
		return log.WithError(utils.ERR_HEADIMG_UPDATE_FAIL)
	}
	user.HeadImg = headImg
	return _self.Update(user)
}
func (_self *UserService) UpdateEmail(id uint64, email string) *utils.Error {
	user, err := _self.SelectOne(id, false)
	if err != nil {
		return log.WithError(utils.ERR_EMAIL_UPDATE_FAIL)
	}
	user.Email = email
	return _self.Update(user)
}
func (_self *UserService) UpdateIntro(id uint64, intro string) *utils.Error {
	if util.Len(intro) < 1 || util.Len(intro) > 30 {
		return log.WithError(utils.ERR_INTRO_VALIDATE_FAIL)
	}
	user, err := _self.SelectOne(id, false)
	if err != nil {
		return log.WithError(utils.ERR_INTRO_UPDATE_FAIL)
	}
	user.Intro = intro
	return _self.Update(user)
}
func (_self *UserService) UpdateNickname(id uint64, nickname string) *utils.Error {
	if util.Len(nickname) < 1 || util.Len(nickname) > 10 {
		return log.WithError(utils.ERR_NICKNAME_VALIDATE_FAIL)
	}
	user, err := _self.SelectOne(id, false)
	if err != nil {
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
func (_self *UserService) UpdateLoginUserKeys(user *entity.User) *utils.Error {
	//没有公钥 创建公私钥
	keys := util.CreateDHKey(conf.Conf.Prime, "02")
	//发起请求修改后台用户信息
	_, err := Post("/api/user/resetPublicKey", &entity.User{PublicKey: keys.PublicKey})
	if err != nil {
		return log.WithError(err)
	}
	user.PublicKey = keys.PublicKey
	user.PrivateKey = keys.PrivateKey
	return nil
}

// Register 用户注册
func (_self *UserService) Register(username, password string) *utils.Error {
	if util.Len(username) < 6 || util.Len(username) > 20 {
		return log.WithError(utils.ERR_USER_USERNAME_LENGTH)
	}
	if util.Len(password) < 6 || util.Len(password) > 20 {
		return log.WithError(utils.ERR_USER_PASSWORD_LENGTH)
	}
	params := &entity.UserReq{
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
	if util.Len(username) < 6 || util.Len(username) > 20 {
		return log.WithError(utils.ERR_USER_USERNAME_LENGTH)
	}
	if util.Len(password) < 6 || util.Len(password) > 20 {
		return log.WithError(utils.ERR_USER_PASSWORD_LENGTH)
	}
	//已经登录过 就不重复登录 但是需要重新获取用户信息 因为可能异地登录导致秘钥更换过
	if conf.GetLoginInfo().Token == "" {
		params := &entity.UserReq{
			Username: username,
			Password: password,
		}
		resultDTO, err := Post("/api/user/login", params)
		if err != nil || resultDTO == nil || resultDTO.Data == nil {
			return log.WithError(err)
		}
		//缓存登录token
		conf.PutToken(resultDTO.Data.(string))
	}
	// 这里标记不需要输入二级密码
	conf.UpdateInputPwd2(-1)
	//获取用户信息
	return _self.LoginInfo()
}
func (_self *UserService) LoginPwd2(pwd2 string) *utils.Error {
	if util.Len(pwd2) < 6 || util.Len(pwd2) > 20 {
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
			e := tx.Commit().Error
			if e != nil {
				return log.WithError(e)
			}
			return nil
		}()
		if err != nil {
			tx.Rollback()
			return err
		}
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
	conf.UpdateInputPwd2(-1)
	//通知前往登录页面
	if Listener != nil {
		Listener.OnLogin()
	}
	return nil
}

// LoginInfo 获取用户信息
func (_self *UserService) LoginInfo() *utils.Error {
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
		resultDTO, err := Post("/api/user/info", nil)
		if err != nil || resultDTO == nil || resultDTO.Data == nil {
			return log.WithError(err)
		}
		//缓存用户信息
		var user entity.User
		e := util.Obj2Obj(resultDTO.Data, &user)
		if e != nil {
			return log.WithError(utils.ERR_GET_USER_INFO_FAIL)
		}

		//服务器公钥是否存在
		if user.PublicKey == "" { //不存在 就重新获取
			err = _self.UpdateLoginUserKeys(&user)
			if err != nil {
				return log.WithError(err)
			}
		}

		//数据库不存在 就添加 这里不做修改
		sysUser, e := _self.repo.Query(&entity.User{Id: user.Id})
		if e != nil {
			log.Error(e)
			return log.WithError(utils.ERR_GET_USER_INFO_FAIL)
		}
		//数据存在--需要把数据库中的私钥封装到登录者中
		if sysUser != nil {
			//公钥存在 但是不一样 或者 数据库里没有私钥
			if sysUser.PublicKey != user.PublicKey || sysUser.PrivateKey == "" {
				err = _self.UpdateLoginUserKeys(&user)
				if err != nil {
					return log.WithError(err)
				}
			} else {
				//都一致  就把数据库的复制给当前的
				user.PrivateKey = sysUser.PrivateKey
			}
		} else {
			//数据库中不存在
			err = _self.UpdateLoginUserKeys(&user)
			if err != nil {
				return log.WithError(err)
			}
		}
		e = _self.Save(&user)
		if e != nil {
			return log.WithError(utils.ERR_GET_USER_INFO_FAIL)
		}
		//覆盖登录文件
		conf.PutLoginInfo(user)
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
