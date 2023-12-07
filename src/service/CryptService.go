package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/util"
	"path/filepath"
	"strings"
	"time"
)

// GetSecret 获取密钥
func GetSecret(target uint64, tp string) (string, *utils.Error) {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return "", log.WithError(utils.ERR_NOT_LOGIN)
	}
	key := tp + "_" + util.Uint642Str(target)
	switch tp {
	case "friend":
		if Keys[key] == "" {
			user, e := QueryUser(target, repository.NewUserRepo())
			if e != nil {
				return "", log.WithError(e)
			}
			if user == nil {
				return "", log.WithError(utils.ERR_ENCRYPT_FAIL)
			}
			secret := util.SharedAESKey(user.PublicKey, conf.GetLoginInfo().User.PrivateKey, conf.Conf.Prime)
			Keys[key] = secret
		}
		break
	case "group":
		break
	}
	return Keys[key], nil
}

// Encrypt 聊天内容加密
func Encrypt(target uint64, tp, content string) (string, *utils.Error) {
	secret, err := GetSecret(target, tp)
	if err != nil {
		return "", err
	}
	data, e := util.EncryptAes(content, secret)
	if e != nil {
		return "", log.WithError(utils.ERR_ENCRYPT_FAIL)
	}
	return data, nil
}

// Decrypt 聊天内容解密
func Decrypt(tp string, target uint64, no, content string) (string, *utils.Error) {
	if content == "" {
		return "", nil
	}
	secret, err := GetSecret(target, tp)
	if err != nil {
		return "", err
	}
	data, err := util.DecryptAes(content, secret)
	//解密失败 直接返回解密失败
	if err != nil {
		return "", log.WithError(err)
	}
	//否则判断是否在聊天中
	if no != "" && conf.Conf.ChatId == target {
		//转换为messageData
		md := &entity.MessageData{}
		e := util.Str2Obj(data, md)
		if e != nil {
			return "", log.WithError(utils.ERR_DECRYPT_FAIL)
		}
		//无需解密文件
		if md.Type < 2 || md.Type > 5 {
			return data, nil
		}
		//需要解密文件
		DecryptFile(target, no, md, secret)
		return util.GetDecryptingMsg(), nil
	}
	return data, nil
}
func DecryptFile(target uint64, no string, md *entity.MessageData, secret string) {
	go func() {
		//延迟2秒
		time.Sleep(2 * time.Second)
		//通过最后一根/获取文件后缀
		paths := strings.Split(md.Content, "/")
		filename := paths[len(paths)-1]
		path := filepath.Join(conf.Base.BaseDir, "configs", filename)
		//先下载文件
		e := util.DownloadFile(md.Content, path)
		if e != nil {
			log.Error(e)
			FileNotify(target, no, util.GetErrMsg())
			return
		}
		//解密文件
		fileData, err := util.DecryptFile(path, secret)
		if err != nil {
			log.Error(err)
			FileNotify(target, no, util.GetErrMsg())
			return
		}
		//保存为临时文件
		tempPath := filepath.Join(conf.Base.BaseDir, "configs", "temp", filename)
		e = util.SaveTempFile(fileData, tempPath)
		if e != nil {
			log.Error(e)
			FileNotify(target, no, util.GetErrMsg())
			return
		}
		okMsg := &entity.MessageData{
			Type:    2,
			Content: tempPath,
		}
		data, e := util.Obj2Str(okMsg)
		if e != nil {
			log.Error(e)
			FileNotify(target, no, util.GetErrMsg())
			return
		}
		FileNotify(target, no, data)
	}()
}
