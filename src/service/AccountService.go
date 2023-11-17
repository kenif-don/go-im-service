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
	if account == nil {
		//没有就从服务器同步
		resultDTO, err := Post("/api/account/selectOne", &entity.Account{UserId: conf.GetLoginInfo().User.Id})
		if err != nil {
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		var a entity.Account
		e = util.Str2Obj(resultDTO.Data.(string), &a)
		if e != nil {
			log.Error(e)
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		e = _self.repo.Save(&a)
		if e != nil {
			log.Error(e)
			return nil, log.WithError(utils.ERR_QUERY_FAIL)
		}
		account = &a
	}
	return account, nil
}
