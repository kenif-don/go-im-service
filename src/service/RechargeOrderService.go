package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"gorm.io/gorm"
)

type IRechargeOrderRepo interface {
	Query(obj *entity.RechargeOrder) (*entity.RechargeOrder, error)
	QueryAll(obj *entity.RechargeOrder) ([]entity.RechargeOrder, error)
	Save(obj *entity.RechargeOrder) error
	Delete(obj *entity.RechargeOrder) error
	BeginTx() *gorm.DB
}
type RechargeOrderService struct {
	repo IRechargeOrderRepo
}
type PayType struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

func NewRechargeOrderService() *RechargeOrderService {
	return &RechargeOrderService{
		repo: repository.NewRechargeOrderRepo(),
	}
}
func QueryRechargeOrder(obj *entity.RechargeOrder, repo IRechargeOrderRepo) (*entity.RechargeOrder, error) {
	return repo.Query(obj)
}
func (_self *RechargeOrderService) GetTypes() *[]PayType {
	//将1-TRC20 2-ERC20封装到*[]PayType中
	return &[]PayType{
		{1, "TRC20"},
		{2, "ERC20"},
	}
}

// AddRechargeOrder 添加充值订单
func (_self *RechargeOrderService) AddRechargeOrder(tp int, value string) (*entity.RechargeOrder, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, utils.ERR_NOT_LOGIN
	}
	reo := &entity.RechargeOrder{
		Value: value,
		Type:  tp,
	}
	// 先添加
	resultDTO, err := Post("/api/recharge-order/selectOne", reo)
	if err != nil {
		return nil, log.WithError(err)
	}
	// 再获取订单进行存储和返回
	id := util.Str2Uint64(resultDTO.Data.(string))
	resultDTO, err = Post("/api/recharge-order/selectOne", map[string]uint64{"id": id})
	if err != nil {
		return nil, log.WithError(err)
	}
	//转换为实体
	rechargeOrder := &entity.RechargeOrder{}
	e := util.Obj2Obj(resultDTO.Data, rechargeOrder)
	if e != nil {
		log.Error(e)
		return nil, utils.ERR_RECHARGE_FAIL
	}
	//保存
	e = _self.repo.Save(rechargeOrder)
	if e != nil {
		log.Error(e)
		return nil, utils.ERR_RECHARGE_FAIL
	}
	return rechargeOrder, nil
}
