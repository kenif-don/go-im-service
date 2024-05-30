package service

import (
	"go-im-service/src/configs/conf"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"go-im-service/src/util"
)

type RechargeOrderService struct {
}
type PayType struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

func NewRechargeOrderService() *RechargeOrderService {
	return &RechargeOrderService{}
}
func (_self *RechargeOrderService) GetTypes() *[]PayType {
	//将1-TRC20 2-ERC20封装到*[]PayType中
	return &[]PayType{
		{1, "TRC20"},
		//{2, "ERC20"},
	}
}

// AddRechargeOrder 添加充值订单
func (_self *RechargeOrderService) AddRechargeOrder(tp int, value string) (string, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return "", utils.ERR_NOT_LOGIN
	}
	if tp == 0 {
		return "", utils.ERR_SELECT_PAY_NETWORK_FAIL
	}
	// 先添加
	resultDTO, err := Post("/api/recharge-order/add", map[string]interface{}{"type": tp, "value": value})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", utils.ERR_RECHARGE_FAIL
	}
	// 再获取订单进行存储和返回
	id := util.Str2Uint64(resultDTO.Data.(string))
	resultDTO, err = Post("/api/recharge-order/selectOne", map[string]uint64{"id": id})
	if err != nil {
		return "", log.WithError(err)
	}
	if resultDTO.Data == nil {
		return "", utils.ERR_RECHARGE_FAIL
	}
	return resultDTO.Data.(string), nil
}
