package service

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/dto"
	"IM-Service/src/util"
	"im-sdk/handler"
	"im-sdk/model"
)

var Listener MessageListener

type MessageListener interface {
	//OnFile 文件解密结果
	OnFile(data []byte)
	//OnReceive 当前聊天接收到消息
	OnReceive(data string)
	//OnDelMsg 如果当前聊天是 对方如果在此时删除了,会触发此通知
	OnDelMsg(data string)
	//OnSendReceive 发送的消息状态 -某消息 1-发送中 2-发送成功、-1-发送失败
	OnSendReceive(data string)
	//OnDoChat 如果客户端停留在首页 如果有新消息进来,都会调用此接口更新最后消息和排序
	OnDoChat(data string)
	//OnFriendApply 好友申请 仅代表有新的好友申请 但是无参
	OnFriendApply()
	//OnLogin 登录失效通知
	OnLogin()
	//OnLoginPwd2 输入二级密码
	OnLoginPwd2()
	//OnDoVoice 收到新消息通知 用来播放语音提示
	OnDoVoice(data string)
	//OnConnectChange 连接状态变化
	OnConnectChange(state string)
}

func SetListener(listener MessageListener) {
	once.Do(func() {
		Listener = listener
		Keys = make(map[string]string)
	})
}

// NotPwdUris 排除输入2级密码的URI
var NotPwdUris = []string{"/api/user/loginPwd2", "/api/offline-bill/selectAll", "/api/offline-bill/dels",
	"/api/group/selectMembers", "/api/friend/selectOne", "/api/friend-apply/selectOne", "/api/user/selectOne",
	"/api/group/selectOne", "/api/agent/selectOne", "/api/config/selectOne"}

func Post(url string, req interface{}) (*dto.ResultDTO, *utils.Error) {
	//排除输入2级密码的URI和需要放行的URI
	if util.IndexOfString(url, NotPwdUris) == -1 && util.IndexOfString(url, conf.Conf.ExUris) == -1 && !ValidatePwd2() {
		log.Debugf("被二级密码拦截的 url: %s", url)
		return nil, log.WithError(utils.ERR_NOT_PWD2_FAIL)
	}
	resultDTO, err := util.Post(url, req)
	if err != nil {
		return nil, log.WithError(err)
	}
	if resultDTO.Code == 401 {
		//退出登录
		err := NewUserService().Logout()
		if err != nil {
			return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
		}
		return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
	}
	if resultDTO.Code == 500 {
		if resultDTO.Msg == "0x99999" {
			return nil, utils.ERR_NET_FAIL
		}
		return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
	}
	return resultDTO, nil
}
func Send(protocol *model.Protocol) *utils.Error {
	mgr := handler.GetClientHandler().GetMessageManager()
	if mgr == nil {
		log.Error("获取消息管理器失败")
		return utils.ERR_SEND_FAIL
	}
	mgr.Send(protocol)
	return nil
}

// ValidatePwd2 判断是否需要输入2级密码 且是否已经输入2级密码 需要且没输 就前往输入2级密码
func ValidatePwd2() bool {
	if conf.GetLoginInfo().User != nil && conf.GetLoginInfo().User.Password2 != "" && conf.GetLoginInfo().InputPwd2 == 1 {
		if Listener != nil {
			//需要输入二级密码 但是没输
			Listener.OnLoginPwd2()
		}
		return false
	}
	return true
}
