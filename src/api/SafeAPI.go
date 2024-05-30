package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/conf"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"
	"go-im-service/src/util"

	"google.golang.org/protobuf/proto"
)

// DeleteOneSafe 删除单个归档记录
func DeleteOneSafe(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.SafeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	//验证密码
	if conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] == "" {
		resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		res, e := proto.Marshal(resp)
		if e != nil {
			return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
		}
		return res
	}
	err := service.NewSafeService().Delete(req.Id)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// InputSafePwd 输入保险箱密码
func InputSafePwd(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.InputSafePwdReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	err := service.NewSafeService().InputSafePwd(req.Pwd)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// AddSafe 归档
func AddSafe(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.SafeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	//验证密码
	if conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] == "" {
		resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		res, e := proto.Marshal(resp)
		if e != nil {
			return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
		}
		return res
	}
	err := service.NewSafeService().Add(req.Content)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(nil, resp)
}

// PagingSafe 分页获取保险箱内容
func PagingSafe(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.SafeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	//验证密码
	if conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] == "" {
		resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		res, e := proto.Marshal(resp)
		if e != nil {
			return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
		}
		return res
	}
	resultDTO, err := service.NewSafeService().Paging(int(req.Page), int(req.PageSize))
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(resultDTO, resp)
}

// SelectOneSafe 获取单个归档
func SelectOneSafe(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.SafeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	//验证密码
	if conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] == "" {
		resp.Code = uint32(api.ResultDTOCode_TO_INPUT_PWD2)
		res, e := proto.Marshal(resp)
		if e != nil {
			return SyncPutErr(utils.ERR_LOGIN_FAIL, resp)
		}
		return res
	}
	resultDTO, err := service.NewSafeService().SelectOne(req.Id)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(resultDTO, resp)
}

// DecrySafeContent 解密保险箱内容
func DecrySafeContent(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.SafeReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	content, err := service.NewSafeService().DecrySafeContent(req.Content)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(content, resp)
}
