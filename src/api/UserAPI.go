package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	"IM-Service/src/service"
	"google.golang.org/protobuf/proto"
)

//type RegisterListener interface {
//	On(data []byte)
//}

// func Register(data []byte, callback RegisterListener)[]byte {
func Register(data []byte) []byte {
	req := &api.RegisterReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(err, resp)
	}
	err := service.NewUserService().Register(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, _ := proto.Marshal(resp)
	return res
}
func Login(data []byte) []byte {
	req := &api.RegisterReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(err, resp)
	}
	// 登录
	err := service.NewUserService().Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
	}
	// 判断释放存在公钥
	if conf.LoginInfo.User.PublicKey == "" {
		//没有公钥 创建公私钥

	}
	// 判断是否存在二级密码

	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, _ := proto.Marshal(resp)
	return res
}
