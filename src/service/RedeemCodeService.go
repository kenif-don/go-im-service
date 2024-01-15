package service

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"strings"
)

type RedeemCodeService struct{}

func NewRedeemCodeService() *RedeemCodeService {
	return &RedeemCodeService{}
}

// Create 钱兑码
func (_self *RedeemCodeService) Create(money string) (string, *utils.Error) {
	if money == "" {
		return "", log.WithError(utils.ERR_PARAM_PARSE)
	}
	resultDTO, err := Post("/api/redeem-code/create", map[string]string{"money": money})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string), nil
}

// Exchange 码兑钱
func (_self *RedeemCodeService) Exchange(code string) *utils.Error {
	if code == "" {
		return log.WithError(utils.ERR_PARAM_PARSE)
	}
	_, err := Post("/api/redeem-code/exchange", map[string]string{"code": strings.TrimSpace(code)})
	if err != nil {
		return log.WithError(err)
	}
	return nil
}

// Paging 分页获取兑换记录
func (_self *RedeemCodeService) Paging(page, pageSize int) (string, *utils.Error) {
	resultDTO, err := Post("/api/redeem-code/paging", map[string]int{"page": page, "pageSize": pageSize})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string), nil
}

// SelectOne 获取兑换记录
func (_self *RedeemCodeService) SelectOne(code string) (string, *utils.Error) {
	if code == "" {
		return "", log.WithError(utils.ERR_PARAM_PARSE)
	}
	resultDTO, err := Post("/api/redeem-code/selectOne", map[string]string{"code": strings.TrimSpace(code)})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", nil
	}
	return resultDTO.Data.(string), nil
}
