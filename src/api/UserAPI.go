package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/im"
	"IM-Service/src/service"
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
	return SyncPutSuccess(user, resp)
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
	return SyncPutSuccess(nil, resp)
}

func Search(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.SearchReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	users, err := service.NewUserService().Search(req.Keyword)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(users, resp)
}
func UpdateSafePwd(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdatePwdReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdatePassword(4, req.Pwd, req.OldPwd, req.NewPwd)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}
func UpdateBurstPwd(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdatePwdReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdatePassword(3, req.Pwd, req.OldPwd, req.NewPwd)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}
func UpdatePwd2(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdatePwdReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdatePassword(2, req.Pwd, req.OldPwd, req.NewPwd)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}
func UpdatePwd(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdatePwdReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdatePassword(1, req.Pwd, req.OldPwd, req.NewPwd)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}
func UpdateHeadImg(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdateHeadImg(conf.GetLoginInfo().User.Id, req.Data)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}
func UpdateEmail(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdateEmail(conf.GetLoginInfo().User.Id, req.Data)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}
func UpdateIntro(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdateIntro(conf.GetLoginInfo().User.Id, req.Data)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}
func UpdateNickname(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UpdateUserReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	err := userService.UpdateNickname(conf.GetLoginInfo().User.Id, req.Data)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
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
	return SyncPutSuccess(conf.GetLoginInfo().User, resp)
}

// AutoLogin 自动登录
func AutoLogin() []byte {
	resp := &api.ResultDTOResp{}
	//存在登录者
	if conf.GetLoginInfo().Token != "" && conf.GetLoginInfo().User != nil && conf.GetLoginInfo().User.Id != 0 {
		//通过info去检查公私钥
		err := service.NewUserService().LoginInfo()
		if err != nil {
			return SyncPutErr(err, resp)
		}
		//判断是否需要输入二级密码
		if conf.GetLoginInfo().InputPwd2 == 1 {
			//需要输入二级密码
			resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
			//需要输入二级密码
			conf.UpdateInputPwd2(1)
		} else {
			resp.Code = uint32(api.ResultDTOCode_SUCCESS)
			conf.UpdateInputPwd2(-1)
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
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	// 需要登录
	userService := service.NewUserService()
	err := userService.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
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
	return SyncPutSuccess(nil, resp)
}
func LoginPwd2(data []byte) []byte {
	req := &api.UserReq{}
	resp := &api.ResultDTOResp{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewUserService().LoginPwd2(req.Password)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

func Register(data []byte) []byte {
	req := &api.UserReq{}
	resp := &api.ResultDTOResp{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewUserServiceNoDB().Register(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}
