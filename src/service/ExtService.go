package service

import (
	"go-im-service/src/configs/conf"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
)

type ExtService struct{}

func NewExtService() *ExtService {
	return &ExtService{}
}

// Get 获取扩展
func (_self *ExtService) Get() (string, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return "", log.WithError(utils.ERR_NOT_LOGIN)
	}
	resultDTO, err := Post("/api/ext/get", nil)
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string) + "?token=" + conf.GetLoginInfo().Token, nil
}
