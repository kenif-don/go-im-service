package service

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/dto"
	"IM-Service/src/util"
	"im-sdk/handler"
	"im-sdk/model"
)

func Post(url string, req interface{}) (*dto.ResultDTO, *utils.Error) {
	resultDTO, e := util.Post(url, req)
	if e != nil {
		return nil, log.WithError(e)
	}
	if resultDTO.Code == 401 {
		//退出登录
		err := NewUserService().Logout()
		if err != nil {
			return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
		}
		return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
	}
	if resultDTO.Code == 500 {
		if resultDTO.Msg == "0x99999" {
			return nil, utils.ERR_NET_FAIL
		}
		return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
	}
	return resultDTO, nil
}
func Send(protocol *model.Protocol) *utils.Error {
	mgr := handler.GetClientHandler().GetMessageManager()
	if mgr == nil {
		return utils.ERR_SEND_FAIL
	}
	mgr.Send(protocol)
	return nil
}
