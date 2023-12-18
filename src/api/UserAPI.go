package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/im"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

// ValidatePwd2 判断是否需要输入2级密码 这个接口会清空2级密码的输入状态
func ValidatePwd2() []byte {
	resp := &api.ResultDTOResp{}
	if conf.GetLoginInfo().User != nil && conf.GetLoginInfo().User.Password2 != "" {
		//需要输入二级密码
		resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		conf.UpdateInputPwd2(1)
	} else {
		//不需要输入二级密码
		resp.Code = uint32(api.ResultDTOCode_SUCCESS)
		conf.UpdateInputPwd2(-1)
	}
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
	}

	return res
}

func SelectOneUser(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UserReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	user, err := service.NewUserService().SelectOne(req.Id, false)
	if err != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	result, e := util.Obj2Str(user)
	if e != nil {
		return SyncPutErr(utils.ERR_QUERY_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = result
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res

}
func Logout() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	err := service.NewUserService().Logout()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}

func Search(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.SearchReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	users, err := userService.Search(req.Keyword)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = users
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}
func UpdateBurstPwd(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdatePwdReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdatePassword(3, req.Pwd, req.OldPwd, req.NewPwd)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}
func UpdatePwd2(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdatePwdReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdatePassword(2, req.Pwd, req.OldPwd, req.NewPwd)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}
func UpdatePwd(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdatePwdReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdatePassword(1, req.Pwd, req.OldPwd, req.NewPwd)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}
func UpdateHeadImg(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateHeadImg(conf.GetLoginInfo().User.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}
func UpdateEmail(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateEmail(conf.GetLoginInfo().User.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}
func UpdateIntro(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateIntro(conf.GetLoginInfo().User.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}
func UpdateNickname(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateNickname(conf.GetLoginInfo().User.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}

	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}

// Info 获取登录者信息
func Info() []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	if conf.GetLoginInfo() == nil || conf.GetLoginInfo().User == nil {
		resp.Code = uint32(api.ResultDTOCode_SUCCESS)
		resp.Msg = "success"
		res, err := proto.Marshal(resp)
		if err != nil {
			return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
		}
		return res
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO_FAIL, resp)
	}
	return res
}

// AutoLogin 自动登录
func AutoLogin() []byte {
	resp := &api.ResultDTOResp{}
	if conf.GetLoginInfo().Token != "" && conf.GetLoginInfo().User != nil && conf.GetLoginInfo().User.Id != 0 {
		//判断是否需要输入二级密码
		if conf.GetLoginInfo().InputPwd2 == 1 {
			//需要输入二级密码
			resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		} else {
			resp.Code = uint32(api.ResultDTOCode_SUCCESS)
		}
		//已经登录--如果已经登录 再链接长连接成功后 会进行一次登录 这里不处理
	} else {
		resp.Code = uint32(api.ResultDTOCode_TO_LOGIN)
	}
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
	}
	return res
}
func Login(data []byte) []byte {
	req := &api.UserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	// 需要登录
	userService := service.NewUserService()
	err := userService.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
	}
	// 判断是否存在公钥
	if conf.LoginInfo.User.PublicKey == "" || conf.LoginInfo.User.PrivateKey == "" {
		//没有公钥 创建公私钥
		keys := util.CreateDHKey(conf.Conf.Prime, "02")
		err = userService.UpdateLoginUserKeys(keys)
		if err != nil {
			return SyncPutErr(err, resp)
		}
		log.Debug("没有私钥，创建私钥")
	}
	//公钥是否和本地一致
	sysUser, err := userService.SelectOne(conf.GetLoginInfo().User.Id, false)
	if err != nil {
		return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
	}
	if sysUser.PublicKey != conf.GetLoginInfo().User.PublicKey || sysUser.PrivateKey != conf.GetLoginInfo().User.PrivateKey {
		//没有公钥 创建公私钥
		keys := util.CreateDHKey(conf.Conf.Prime, "02")
		err = userService.UpdateLoginUserKeys(keys)
		if err != nil {
			return SyncPutErr(err, resp)
		}
		log.Debug("有私钥，但是私钥不一致，更换私钥")
	}
	// 判断是否存在二级密码
	if conf.LoginInfo.User.Password2 != "" {
		//存在2级密码 跳转到输入二级密码页面
		resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		//需要输入二级密码
		conf.UpdateInputPwd2(1)
	} else {
		//不需要输入二级密码
		conf.UpdateInputPwd2(-1)
		resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	}
	//登录IM
	err = im.LoginIm()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
	}
	return res
}
func LoginPwd2(data []byte) []byte {
	req := &api.UserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewUserService().LoginPwd2(req.Password)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
	}
	return res
}

func Register(data []byte) []byte {
	req := &api.UserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewUserServiceNoDB().Register(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
	}
	return res
}
