package service

import (
	"go-im-service/src/configs/conf"
	_const "go-im-service/src/configs/const"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"go-im-service/src/entity"
	"go-im-service/src/repository"
	"go-im-service/src/util"

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
func (_self *AccountService) Transfer(tp, remark, amount, password, no string, gId, he uint64) *utils.Error {
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
		"no":       no,
	}
	_, err := Post("/api/account/transfer", req)
	if err != nil {
		return log.WithError(err)
	}
	req["password"] = nil
	req["no"] = nil
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
		return NewMessageService().SendMsg(tp, he, no, md)
	}
	return NewMessageService().SendMsg(tp, gId, no, md)
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
	if resultDTO.Data == nil {
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
