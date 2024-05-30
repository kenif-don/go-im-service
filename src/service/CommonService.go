package service

import (
	"go-im-service/src/configs/conf"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"go-im-service/src/util"
)

// SelectConfig 获取平台配置
func SelectConfig() (string, *utils.Error) {
	resultDTO, err := Post("/api/config/selectOne", map[string]interface{}{})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO != nil {
		res, e := util.Obj2Str(resultDTO.Data)
		if e != nil {
			log.Error(e)
			return "", log.WithError(utils.ERR_QUERY_FAIL)
		}
		return res, nil
	}
	return "", nil
}

// SetLanguage 设置语言zh-CN/en
func SetLanguage(language string) {
	conf.Conf.Language = language
}
func GetAgent() (string, *utils.Error) {
	resultDTO, err := Post("/api/agent/selectOne", nil)
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO != nil {
		res, e := util.Obj2Str(resultDTO.Data)
		if e != nil {
			log.Error(e)
			return "", log.WithError(utils.ERR_QUERY_FAIL)
		}
		return res, nil
	}
	return "", nil
}
