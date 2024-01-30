package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/util"
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
