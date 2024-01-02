package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/util"
	"bytes"
	"github.com/go-audio/wav"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strconv"
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
			user, e := NewUserService().SelectOne(target, false)
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
		if Keys[key] == "" {
			group, e := NewGroupService().SelectOne(target, false)
			if e != nil {
				return "", log.WithError(e)
			}
			if group == nil {
				return "", log.WithError(utils.ERR_ENCRYPT_FAIL)
			}
			//加密群
			if group.Type == 2 {
				secret := util.MD5("group_" + group.Password)
				Keys[key] = secret
			} else {
				secret := util.MD5("group_" + strconv.FormatUint(group.Id, 10))
				Keys[key] = secret
			}
		}
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
		return "", log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	//字符串转换为消息体
	md := &entity.MessageData{}
	e := util.Str2Obj(content, md)
	if e != nil {
		return "", log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	//系统消息  无需解密
	if md.Type == 9 {
		return content, nil
	}
	//先解密消息体
	secret, err := GetSecret(target, tp)
	if err != nil {
		return "", err
	}
	md.Content, err = util.DecryptAes(md.Content, secret)
	data, e := util.Obj2Str(md)
	if e != nil {
		return "", log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	//解密失败 直接返回解密失败
	if err != nil {
		return "", log.WithError(err)
	}
	if md.Type < 2 || md.Type > 5 {
		return data, nil
	}
	//否则判断是否在聊天中
	if no != "" && conf.Conf.ChatId == target {
		return util.GetDecryptingMsg(md), nil
	}
	return data, nil
}
func DecryptFile(tp string, target uint64, no string) *utils.Error {
	log.Debugf("11111111111111")
	//如果当前聊天不是正在聊天的 就不解密了
	if conf.Conf.ChatId != target {
		return nil
	}
	log.Debugf("222222222222")
	message, err := NewMessageService().SelectOne(&entity.Message{No: no})
	if err != nil || message == nil {
		log.Error(err)
		return log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	log.Debugf("3333333333333333")
	var md = &entity.MessageData{}
	e := util.Str2Obj(message.Data, md)
	if e != nil {
		return log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	chat, e := NewChatService().repo.Query(&entity.Chat{Type: tp, TargetId: target})
	if e != nil || chat == nil {
		log.Error(e)
		return log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	log.Debugf("444444444444444444444")
	//不是PC 无需解密文件
	if md.Type < 2 || md.Type > 5 {
		return nil
	}
	if md.Type == 5 && conf.Base.DeviceType != conf.PC {
		return nil
	}
	log.Debugf("55555555555555555555555")
	//如果是文件消息 需要解密
	secret, err := GetSecret(chat.TargetId, chat.Type)
	if err != nil {
		return err
	}
	go func() {
		log.Debugf("66666666666666666")
		md.Content, err = util.DecryptAes(md.Content, secret)
		if err != nil {
			log.Error(err)
			FileNotify(chat.TargetId, no, util.GetErrMsg(md.Type))
			return
		}
		//通过最后一根/获取文件后缀
		paths := strings.Split(md.Content, "/")
		filename := paths[len(paths)-1]
		path := filepath.Join(conf.Base.BaseDir, "configs", filename)
		log.Debugf("777777777777777777")
		//先下载文件
		e := util.DownloadFile(md.Content, path)
		if e != nil {
			log.Error(e)
			FileNotify(chat.TargetId, no, util.GetErrMsg(md.Type))
			return
		}
		log.Debugf("88888888888888888")
		//解密文件
		fileData, err := util.DecryptFile(path, secret)
		if err != nil {
			log.Error(err)
			FileNotify(chat.TargetId, no, util.GetErrMsg(md.Type))
			return
		}
		log.Debugf("9999999999999999999")
		//保存为临时文件
		tempPath := filepath.Join(conf.Base.BaseDir, "configs", "temp", filename)
		e = util.SaveTempFile(fileData, tempPath)
		if e != nil {
			log.Error(e)
			FileNotify(chat.TargetId, no, util.GetErrMsg(md.Type))
			return
		}
		log.Debugf("0000000000000000000000")
		//组装成各种类型
		okMsg, err := coverMessageData(md, fileData, tempPath)
		if err != nil {
			log.Error(err)
			FileNotify(chat.TargetId, no, util.GetErrMsg(md.Type))
			return
		}
		log.Debugf("aaaaaaaaaaaaaaaaaaaaa")
		data, e := util.Obj2Str(okMsg)
		if e != nil {
			log.Error(e)
			FileNotify(chat.TargetId, no, util.GetErrMsg(md.Type))
			return
		}
		log.Debugf("bbbbbbbbbbbbbbbbbbbb")
		FileNotify(chat.TargetId, no, data)
		log.Debugf("ccccccccccccccccc")
	}()
	return nil
}

func coverMessageData(md *entity.MessageData, data []byte, path string) (*entity.MessageData, *utils.Error) {
	switch md.Type {
	case 2: //图片
		//获取文件后缀
		endWidth, err := util.GetFileType(path, data)
		if err != nil {
			log.Error(err)
			return nil, log.WithError(utils.ERR_DECRYPT_FAIL)
		}
		c, e := DecodeImageWidthHeight(data, endWidth)
		if e != nil {
			log.Error(e)
			return nil, log.WithError(utils.ERR_DECRYPT_FAIL)
		}
		return &entity.MessageData{
			Type:    md.Type,
			Content: path,
			Status:  2,
			Width:   c.Width,
			Height:  c.Height,
		}, nil
	case 3: //语音
		decoder := wav.NewDecoder(bytes.NewReader(data))
		if decoder == nil {
			return nil, log.WithError(utils.ERR_DECRYPT_FAIL)
		}
		duration, e := decoder.Duration()
		if e != nil {
			return nil, log.WithError(utils.ERR_DECRYPT_FAIL)
		}
		return &entity.MessageData{
			Type:     md.Type,
			Content:  path,
			Status:   2,
			Duration: int(duration.Seconds()),
		}, nil
	case 4: //视频
		return &entity.MessageData{
			Type:    md.Type,
			Content: path,
			Status:  2,
		}, nil
	case 5: //文件
		return &entity.MessageData{
			Type:    md.Type,
			Content: path,
			Size:    len(data),
			Status:  2,
		}, nil
	}
	return nil, nil
}

// DecodeImageWidthHeight 解析图片的宽高信息
func DecodeImageWidthHeight(imgBytes []byte, fileType string) (image.Config, error) {
	switch strings.ToLower(fileType) {
	case ".jpg", ".jpeg":
		return jpeg.DecodeConfig(bytes.NewReader(imgBytes))
	//case ".webp":
	//	imgConf, err = webp.DecodeConfig(bytes.NewReader(imgBytes))
	case ".png":
		return png.DecodeConfig(bytes.NewReader(imgBytes))
	//case ".tif", ".tiff":
	//	imgConf, err = tiff.DecodeConfig(bytes.NewReader(imgBytes))
	case ".gif":
		return gif.DecodeConfig(bytes.NewReader(imgBytes))
		//case "bmp":
		//	imgConf, err = bmp.DecodeConfig(bytes.NewReader(imgBytes))
	}
	return image.Config{}, nil
}
