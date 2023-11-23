package im

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"im-sdk/client"
	"im-sdk/handler"
	"im-sdk/model"
	"strconv"
	"time"
)

type LogicProcess struct{}

var process = &LogicProcess{}

func GetLogicProcess() *LogicProcess {
	return process
}
func StartIM() {
	go func() {
		for {
			//启动长连接
			conf.Conf.Connected = true
			conf.Conf.Client = client.New(conf.Base.WsHost)
			e := conf.Conf.Client.Startup(GetLogicProcess())
			if e != nil {
				_ = log.WithError(e, "启动长连接失败，准备重启")
				conf.Conf.Connected = false
			}
			if conf.Conf.Connected {
				return
			}
			time.Sleep(time.Second * 2)
		}
	}()
}
func (_self *LogicProcess) Connected() {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return
	}
	//获取登录者 组装登录IM请求参数
	loginInfo := &model.LoginInfo{
		Id:     strconv.FormatUint(conf.GetLoginInfo().User.Id, 10),
		Device: conf.Base.DeviceType,
		Token:  conf.GetLoginInfo().Token,
	}
	handler.GetClientHandler().GetMessageManager().SendLogin(loginInfo)
	if service.Listener != nil {
		service.Listener.OnConnectChange("1")
	}
}

// SendOk qos中的消息发送成功 服务器成功返回
func (_self *LogicProcess) SendOk(protocol *model.Protocol) {
	messageService := service.NewMessageService()
	messageService.UpdateReaded(protocol, 2)
	if service.Listener != nil && (protocol.Type == 1 || protocol.Type == 8) {
		res, e := util.Obj2Str(map[string]interface{}{"no": protocol.No, "send": 2})
		if e != nil {
			log.Error(e)
			return
		}
		//消息状态通知
		service.Listener.OnSendReceive(res)
		var message = &entity.Message{}
		e = util.Str2Obj(protocol.Data.(string), message)
		if e != nil {
			log.Error(e)
			return
		}
		//判断是否存在聊天--这里是自己发的 所以聊天的目标是消息的目标
		chat, e := service.QueryChat(message.Type, message.TargetId, repository.NewChatRepo())
		if e != nil {
			log.Error(e)
			return
		}
		if chat == nil {
			chat, e = service.NewChatService().CoverChat(message.Type, util.Str2Uint64(protocol.From))
			if e != nil {
				log.Error(e)
				return
			}
		}
		// 通知聊天列表更新
		err := service.NewChatService().ChatNotify(chat)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

// SendOkCallback 发送成功的回调
// 仅仅是发出去了 如果是Qos消息 此时还未收到服务器反馈
// SendOk代表发出Qos消息并接收到了服务器反馈
func (_self *LogicProcess) SendOkCallback(protocol *model.Protocol) {

}

// SendFailedCallback 发送失败的回调
func (_self *LogicProcess) SendFailedCallback(protocol *model.Protocol) {
	messageService := service.NewMessageService()
	messageService.UpdateReaded(protocol, -1)
}

// LoginOk 登录成功的回调
func (_self *LogicProcess) LoginOk(protocol *model.Protocol) {
	conf.Conf.LoginIM = true
	conf.DiffTime = int(util.Str2Uint64(protocol.Data.(string)) - util.CurrentTime())
	log.Debugf("登录成功！时差:%v", conf.DiffTime)
	//获取一次通讯录
	_, e := service.NewFriendService().SelectAll()
	//获取离线消息
	e = service.NewMessageService().GetOfflineMessage()
	if e != nil {
		log.Error(e)
	}
}

// LoginFail 登录失败的回调
func (_self *LogicProcess) LoginFail(protocol *model.Protocol) {
	log.Debugf("登录失败！:%v", protocol)
}

// Logout 客户端正常退出
func (_self *LogicProcess) Logout() {
	//进行重连
	conf.Conf.LoginIM = false
	go func() {
		err := conf.Conf.Client.Reconnect()
		if err != nil {
			log.Error(err)
		}
	}()
}

// ReceivedMessage 接收到消息
func (_self *LogicProcess) ReceivedMessage(protocol *model.Protocol) {
	log.Debugf("接收到服务器IM消息:%v", protocol)
	err := service.NewMessageService().Handler(protocol)
	if err != nil {
		log.Errorf("解析服务器IM消息失败:%v", err)
	}
}
func (_self *LogicProcess) Exception(msg string) {
	if service.Listener != nil {
		service.Listener.OnConnectChange("0")
	}
	log.Errorf("exception:%v", msg)
	conf.Conf.Connected = false
	log.Debug("服务器断开连接,进行重连")
	go func() {
		for {
			err := conf.Conf.Client.Reconnect()
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(5 * time.Second)
		}
	}()
}
