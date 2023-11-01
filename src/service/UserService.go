package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"context"
	"gorm.io/gorm"
)

type IUserRepo interface {
	Query(obj *entity.User) (*entity.User, error)
	QueryAll(obj *entity.User) (*[]entity.User, error)
	Update(obj *entity.User) error
	Create(obj *entity.User) error
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
		repo: repository.NewUserRepo(context.Background()),
	}
}
func QueryUser(id uint64, repo IUserRepo) (*entity.User, error) {
	return repo.Query(&entity.User{Id: id})
}
func (_self *UserService) Create(obj *entity.User) error {
	return _self.repo.Create(obj)
}

// UpdateLoginUserKeys 修改公私钥
func (_self *UserService) UpdateLoginUserKeys(keys EncryptKeys) error {
	obj, err := QueryUser(conf.LoginInfo.User.Id, _self.repo)
	if err != nil {
		return log.WithError(err)
	}
	if obj == nil {
		return log.WithError(utils.ERR_USER_NOT_EXIST)
	}
	obj.PublicKey = keys.PublicKey
	obj.PrivateKey = keys.PrivateKey
	tx := _self.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return log.WithError(err)
	}
	//发起请求修改后台用户信息
	_, err = util.Post("/api/user/resetPublicKey", &entity.User{PublicKey: keys.PublicKey})
	if err != nil {
		return log.WithError(err)
	}
	//修改本地信息
	err = _self.repo.Update(obj)
	if err != nil {
		return log.WithError(err)
	}
	//修改登录者
	conf.PutLoginInfo(*obj)
	return tx.Commit().Error
}

// Register 用户注册
func (_self *UserService) Register(username, password string) (err error) {
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
	_, err = util.Post("/api/user/register", params)
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

func (_self *UserService) Login(username, password string) (err error) {
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
	resultDTO, err := util.Post("/api/user/login", params)
	if err != nil {
		return log.WithError(err)
	}
	if resultDTO.Data == nil {
		return utils.NewError(500, "login fail", "登录失败")
	}
	//缓存登录token
	conf.PutToken(resultDTO.Data.(string))
	// 不知道是否需要输入
	conf.UpdateInputPwd2(1)
	//获取用户信息
	err = _self.LoginInfo()
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

// LoginInfo 获取用户信息
func (_self *UserService) LoginInfo() (err error) {
	resultDTO, err := util.Post("/api/user/info", nil)
	if err != nil {
		return log.WithError(err)
	}
	if resultDTO.Data == nil {
		return log.WithError(utils.NewError(500, "login info get fail", "登录信息获取失败"))
	}
	//缓存用户信息
	var user entity.User
	err = util.Map2Obj(resultDTO.Data, &user)
	if err != nil {
		return log.WithError(utils.NewError(500, "login info get fail", "登录信息获取失败"))
	}
	//存到文件
	conf.PutLoginInfo(user)
	//数据库不存在 就添加 这里不做修改
	sysUser, err := QueryUser(user.Id, _self.repo)
	if err != nil {
		return log.WithError(err)
	}
	if sysUser != nil {
		return nil
	}
	err = _self.Create(&user)
	if err != nil {
		return log.WithError(err)
	}
	return nil
}
