package service

import (
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"go-im-service/src/util"
)

type VersionService struct {
}

func NewVersionService() *VersionService {
	return &VersionService{}
}
func (_self *VersionService) GetVersion(versionCode, tp int32) (string, *utils.Error) {
	resultDTO, err := Post("/public/version/select", map[string]int32{"versionCode": versionCode, "type": tp})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", nil
	}
	res, e := util.Obj2Str(resultDTO.Data)
	if e != nil {
		return "", log.WithError(utils.ERR_NET_FAIL)
	}
	return res, nil
}
