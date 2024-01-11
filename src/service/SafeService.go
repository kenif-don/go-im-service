package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/util"
	"strings"
)

type SafeService struct{}

func NewSafeService() *SafeService {
	return &SafeService{}
}

// InputSafePwd 修改安全密码
func (_self *SafeService) InputSafePwd(pwd string) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return log.WithError(utils.ERR_NOT_LOGIN)
	}
	if pwd == "" || conf.GetLoginInfo().User.SafePassword == "" {
		return log.WithError(utils.ERR_INPUT_SAFE_PASSWORD)
	}
	secret := util.MD5(pwd)
	if strings.ToUpper(secret) != conf.GetLoginInfo().User.SafePassword {
		return log.WithError(utils.ERR_SAFE_PASSWORD)
	}
	conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] = util.MD5("safe_" + pwd)
	return nil
}

// Add 归档
func (_self *SafeService) Add(content string) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return log.WithError(utils.ERR_NOT_LOGIN)
	}
	//验证密码
	if conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] == "" {
		return log.WithError(utils.ERR_INPUT_SAFE_PASSWORD)
	}
	encryptStr, err := util.EncryptAes(content, conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)])
	if err != nil {
		return log.WithError(err)
	}
	_, err = Post("/api/safe/add", map[string]string{"content": encryptStr})
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

// Paging 分页获取归档
func (_self *SafeService) Paging(page, pageSize int) ([]entity.Safe, *utils.Error) {
	//验证密码
	if conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] == "" {
		return nil, log.WithError(utils.ERR_INPUT_SAFE_PASSWORD)
	}
	resultDTO, err := Post("/api/safe/paging", map[string]int{"page": page, "pageSize": pageSize})
	if err != nil {
		return nil, log.WithError(err)
	}
	var safes []entity.Safe
	e := util.Str2Obj(resultDTO.Data.(string), &safes)
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//循环解密
	for i := 0; i < len(safes); i++ {
		data, err := util.DecryptAes(safes[i].Content, conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)])
		if err != nil {
			safes[i].Content = "解密失败"
		} else {
			safes[i].Content = data
		}
	}
	return safes, nil
}

// SelectOne 获取单个归档
func (_self *SafeService) SelectOne(id uint64) (*entity.Safe, *utils.Error) {
	//验证密码
	if conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)] == "" {
		return nil, log.WithError(utils.ERR_INPUT_SAFE_PASSWORD)
	}
	resultDTO, err := Post("/api/safe/selectOne", map[string]uint64{"id": id})
	if err != nil {
		return nil, log.WithError(err)
	}
	var safe entity.Safe
	e := util.Str2Obj(resultDTO.Data.(string), &safe)
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//解密
	data, err := util.DecryptAes(safe.Content, conf.Conf.Pwds["safe_"+util.Uint642Str(conf.GetLoginInfo().User.Id)])
	if err != nil {
		safe.Content = "解密失败"
	} else {
		safe.Content = data
	}
	return &safe, nil
}
