package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

//	type RegisterListener interface {
//		On(data []byte)
//	}
//
// Info 获取登录者信息
func Info() []byte {
	resp := &api.ResultDTOResp{}
	user, err := util.Obj2Str(conf.GetLoginInfo().User)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	resp.Data = user
	res, err := proto.Marshal(resp)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return res
}
func Login(data []byte) []byte {
	req := &api.RegisterReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(err, resp)
	}
	// 判断是否有登录者
	if conf.GetLoginInfo().Token != "" {
		//判断是否需要输入二级密码
		if conf.GetLoginInfo().InputPwd2 == 2 {
			//需要输入二级密码
			resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		} else {
			resp.Code = uint32(api.ResultDTOCode_SUCCESS)
		}
		resp.Msg = "success"
		res, _ := proto.Marshal(resp)
		return res
	}
	// 需要登录
	userService := service.NewUserService()
	err := userService.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
	}
	// 判断是否存在公钥
	if conf.LoginInfo.User.PublicKey == "" {
		//没有公钥 创建公私钥
		keys := service.CreateDHKey("262074f1e0e19618f0d2af786779d6ad9e814b", "02")
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
	res, _ := proto.Marshal(resp)
	return res
}

// func Register(data []byte, callback RegisterListener)[]byte {
func Register(data []byte) []byte {
	req := &api.RegisterReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(err, resp)
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
