package service

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
)

type VersionService struct {
}

func NewVersionService() *VersionService {
	return &VersionService{}
}
func (_self *VersionService) GetVersion(versionCode, tp int32) (string, *utils.Error) {
	resultDTO, err := Post("/api/version/select", map[string]int32{"versionCode": versionCode, "type": tp})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string), nil
}
