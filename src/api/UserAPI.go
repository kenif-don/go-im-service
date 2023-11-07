package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/repository"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
	"im-sdk/handler"
	"im-sdk/model"
	"strconv"
)

func Logout() []byte {
	resp := &api.ResultDTOResp{}
	conf.ClearLoginInfo()
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
func Search(data []byte) []byte {
	req := &api.SearchReq{}
	resp := &api.ResultDTOResp{}
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
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
func UpdateHeadImg(data []byte) []byte {
	req := &api.UpdateUserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateHeadImg(req.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
func UpdateEmail(data []byte) []byte {
	req := &api.UpdateUserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateEmail(req.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
func UpdateIntro(data []byte) []byte {
	req := &api.UpdateUserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateIntro(req.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
func UpdateNickname(data []byte) []byte {
	req := &api.UpdateUserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	userService := service.NewUserService()
	e := userService.UpdateNickname(req.Id, req.Data)
	if e != nil {
		return SyncPutErr(e, resp)
	}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}

	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res

}

// Info 获取登录者信息
func Info() []byte {
	resp := &api.ResultDTOResp{}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Body = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(utils.ERR_GET_USER_INFO, resp)
	}
	return res
}
func Login(data []byte) []byte {
	req := &api.UserReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	//用户没有传账号密码 直接通过缓存登录
	if req.Username == "" && conf.GetLoginInfo().Token != "" {
		//判断是否需要输入二级密码
		if conf.GetLoginInfo().InputPwd2 == 2 {
			//需要输入二级密码
			resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		} else {
			resp.Code = uint32(api.ResultDTOCode_SUCCESS)
		}
		resp.Msg = "success"
		res, e := proto.Marshal(resp)
		if e != nil {
			return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
		}
		//登录IM
		loginIM()
		return res
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
	}
	//公钥是否和本地一致
	sysUser, e := service.QueryUser(conf.GetLoginInfo().User.Id, repository.NewUserRepo())
	if e != nil {
		return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
	}
	if sysUser.PublicKey != conf.GetLoginInfo().User.PublicKey || sysUser.PrivateKey != conf.GetLoginInfo().User.PrivateKey {
		//没有公钥 创建公私钥
		keys := util.CreateDHKey(conf.Conf.Prime, "02")
		err = userService.UpdateLoginUserKeys(keys)
		if err != nil {
			return SyncPutErr(err, resp)
		}
	}
	// 判断是否存在二级密码
	if conf.LoginInfo.User.Password2 != "" {
		//存在2级密码 跳转到输入二级密码页面
		resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		//需要输入二级密码
		conf.UpdateInputPwd2(2)
	} else {
		//不需要输入二级密码
		conf.UpdateInputPwd2(-1)
		resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	}
	resp.Msg = "success"
	//登录IM
	loginIM()
	res, _ := proto.Marshal(resp)

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
	res, _ := proto.Marshal(resp)
	return res
}

// loginIM 如果已经存在登录者 直接登录 否则走完整登录流程 然后再登录IM
func loginIM() {
	//获取登录者 组装登录IM请求参数
	loginInfo := &model.LoginInfo{
		Id:     strconv.FormatUint(conf.GetLoginInfo().User.Id, 10),
		Device: conf.Base.DeviceType,
		Token:  conf.GetLoginInfo().Token,
	}
	handler.GetClientHandler().GetMessageManager().SendLogin(loginInfo)
}
