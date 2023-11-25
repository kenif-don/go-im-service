package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
)

type WithdrawalService struct{}

func NewWithdrawalService() *WithdrawalService {
	return &WithdrawalService{}
}

// GetWithdrawalFee 获取提现手续费
func (_self *WithdrawalService) GetWithdrawalFee() (string, *utils.Error) {
	resultDTO, err := Post("/api/config/selectOne", map[string]string{})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", utils.ERR_GET_WITHDRAWAL_FEE_FAIL
	}
	return resultDTO.Data.(string), nil
}

// AddWithdrawal 添加提现
func (_self *WithdrawalService) AddWithdrawal(money float64, address string) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return utils.ERR_NOT_LOGIN
	}
	if money <= 0 {
		return utils.ERR_INPUT_WITHDRAWAL_MONEY_FAIL
	}
	if address == "" {
		return utils.ERR_INPUT_WITHDRAWAL_WALLET_FAIL
	}
	// 先添加
	_, err := Post("/api/withdrawal/add", map[string]interface{}{"money": money, "address": address})
	if err != nil {
		return log.WithError(err)
	}
	return nil
}
