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
func Decrypt(target uint64, tp, no, content string) (string, *utils.Error) {
	if content == "" {
		return "", nil
	}
	secret, err := GetSecret(target, tp)
	if err != nil {
		return "", err
	}
	data, err := util.DecryptAes(content, secret)
	if err != nil {
		msg := &entity.MessageData{
			Type:    1,
			Content: "解密失败",
		}
		d, e := util.Obj2Str(msg)
		data = d
		if e != nil {
			return "", log.WithError(utils.ERR_DECRYPT_FAIL)
		}
	}
	//有消息ID才需要解密动作
	if no != "" && conf.Conf.ChatId == target {
		//解密文件
		data, err = DecryptFile(no, data, secret)
		if err != nil {
			msg := &entity.MessageData{
				Type:    1,
				Content: "文件解密失败",
			}
			d, e := util.Obj2Str(msg)
			if e != nil {
				return "", log.WithError(utils.ERR_DECRYPT_FAIL)
			}
			data = d
		}
	}
	return data, nil
}
func DecryptFile(no, data, secret string) (string, *utils.Error) {
	//转换为messageData
	md := &entity.MessageData{}
	e := util.Str2Obj(data, md)
	if e != nil {
		return "", utils.ERR_DECRYPT_FAIL
	}
	if md.Type < 2 && md.Type > 5 {
		return data, nil
	}
	//如果是文件 进入解密
	go func() {
		//先下载文件
		//通过最后一根/获取文件后缀
		paths := strings.Split(md.Content, "/")
		filename := paths[len(paths)-1]
		path := filepath.Join(conf.Base.BaseDir, "configs", filename)
		e := util.DownloadFile(md.Content, path)
		if e != nil {
			log.Error(e)
			//通知文件解密失败
			FileNotify(no, "", -1, int32(md.Type), "", nil)
			return
		}
		//解密文件
		fileData, err := util.DecryptFile(path, secret)
		if err != nil {
			log.Error(err)
			//通知文件解密失败
			FileNotify(no, "", -1, int32(md.Type), "", nil)
			return
		}
		//保存为临时文件
		tempPath := filepath.Join(conf.Base.BaseDir, "configs", "temp", filename)
		e = util.SaveTempFile(fileData, tempPath)
		if e != nil {
			log.Error(e)
			//通知文件解密失败
			FileNotify(no, tempPath, -1, int32(md.Type), "", nil)
			return
		}
		//通知文件解密成功
		FileNotify(no, tempPath, 1, int32(md.Type), "", nil)
	}()
	msg := &entity.MessageData{
		Type:    1,
		Content: "文件解密中",
	}
	data, e = util.Obj2Str(msg)
	if e != nil {
		return "", log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	return data, nil
}
