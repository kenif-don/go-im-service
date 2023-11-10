package im

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/log"
	"IM-Service/src/service"
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
}

// SendOk qos中的消息发送成功 服务器成功返回
func (_self *LogicProcess) SendOk(protocol *model.Protocol) {
	messageService := service.NewMessageService()
	messageService.UpdateReaded(protocol, 2)
	if service.Listener != nil {
		service.Listener.OnSendReceive([]byte("123"))
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
	log.Debugf("登录成功！:%v", protocol)
	//获取离线消息
	e := service.NewMessageService().GetOfflineMessage()
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
	log.Errorf("exception:%v", msg)
	if msg == "unexpected EOF" || msg == "ws closed: 1000 Bye" {
		conf.Conf.Connected = false
		log.Debug("服务器断开连接,进行重连")
		go func() {
			err := conf.Conf.Client.Reconnect()
			if err != nil {
				log.Error(err)
			}
		}()
	}
}
