package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
)

type CashFlowService struct {
}

func NewCashFlowService() *CashFlowService {
	return &CashFlowService{}
}
func (_self *CashFlowService) Paging(page int) (string, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return "", utils.ERR_NOT_LOGIN
	}
	resultDTO, err := Post("/api/cash-flow/paging", map[string]interface{}{"page": page})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", log.WithError(utils.ERR_QUERY_FAIL)
	}
	return resultDTO.Data.(string), nil
}
