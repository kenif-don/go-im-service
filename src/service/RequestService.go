package service

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/dto"
	"IM-Service/src/util"
)

func Post(url string, req interface{}) (*dto.ResultDTO, *utils.Error) {
	resultDTO, e := util.Post(url, req)
	if e != nil {
		return nil, log.WithError(e)
	}
	if resultDTO.Code == 401 {
		if Listener != nil {
			Listener.OnLogin()
		}
		return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
	}
	return resultDTO, nil
}
