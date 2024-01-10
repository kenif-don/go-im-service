package service

import (
	"IM-Service/src/configs/conf"
	_const "IM-Service/src/configs/const"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"gorm.io/gorm"
)

type IAccountRepo interface {
	Query(obj *entity.Account) (*entity.Account, error)
	QueryAll(obj *entity.Account) ([]entity.Account, error)
	Save(obj *entity.Account) error
	Delete(obj *entity.Account) error
	BeginTx() *gorm.DB
}
type AccountService struct {
	repo IAccountRepo
}

func NewAccountService() *AccountService {
	return &AccountService{
		repo: repository.NewAccountRepo(),
	}
}

// Transfer 转账并发消息
func (_self *AccountService) Transfer(tp, remark, amount, password string, gId, he uint64) *utils.Error {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return utils.ERR_NOT_LOGIN
	}
	req := map[string]interface{}{
		"type":     tp,
		"remark":   remark,
		"gId":      gId,
		"he":       he,
		"amount":   amount,
		"password": password,
	}
	_, err := Post("/api/account/transfer", req)
	if err != nil {
		return log.WithError(err)
	}
	req["password"] = nil
	reqStr, e := util.Obj2Str(req)
	if e != nil {
		log.Error(e)
		return log.WithError(utils.ERR_TRANSFER_FAIL)
	}
	md := &entity.MessageData{
		Type:    _const.MSG_TRANSFER,
		Content: reqStr,
	}
	// 发送消息
	if gId == 0 {
		return NewMessageService().SendMsg(tp, he, util.GetUUID(), md)
	}
	return NewMessageService().SendMsg(tp, gId, util.GetUUID(), md)
}

// SelectOneAccount 获取登录者账户数据 没有就从服务器同步
func (_self *AccountService) SelectOneAccount(flush bool) (*entity.Account, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil, utils.ERR_NOT_LOGIN
	}
	account, e := _self.repo.Query(&entity.Account{UserId: conf.GetLoginInfo().User.Id})
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	//有数据 并且不使用远程就直接返回
	if account != nil && !flush {
		return account, nil
	}
	//没有就从服务器同步
	resultDTO, err := Post("/api/account/selectOne", &entity.Account{UserId: conf.GetLoginInfo().User.Id})
	if err != nil {
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	a := &entity.Account{}
	e = util.Str2Obj(resultDTO.Data.(string), a)
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	e = _self.repo.Save(a)
	if e != nil {
		log.Error(e)
		return nil, log.WithError(utils.ERR_QUERY_FAIL)
	}
	return a, nil
}
