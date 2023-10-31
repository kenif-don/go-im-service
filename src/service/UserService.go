package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/util"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
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

// Login 用户登录
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
	conf.LoginInfo.Token = resultDTO.Data.(string)
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
	conf.LoginInfo.User = &user
	//数据库不存在 就添加 这里不做修改

	return nil
}
