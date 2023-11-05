package im

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"im-sdk/client"
	"im-sdk/handler"
	"im-sdk/model"
)

type LogicProcess struct{}

var process = &LogicProcess{}

func GetLogicProcess() *LogicProcess {
	return process
}
func StartIM() {
	//启动长连接
	conf.Conf.Connected = true
	c := client.New(conf.Base.WsHost)
	e := c.Startup(GetLogicProcess())
	if e != nil {
		_ = log.WithError(e, "启动长连接失败")
		conf.Conf.Connected = false
	}
	go func() {
		mgr := handler.GetClientHandler().GetMessageManager()
		//开启心跳
		mgr.StartupHeartbeat()
		//开启Qos
		mgr.StartupQos()
	}()
}

// SendOkCallback 发送成功的回调
// 仅仅是发出去了 如果是Qos消息 此时还未收到服务器反馈
// SendOk代表发出Qos消息并接收到了服务器反馈
func (_self *LogicProcess) SendOkCallback(protocol *model.Protocol) {

}

// SendFailedCallback 发送失败的回调
func (_self *LogicProcess) SendFailedCallback(protocol *model.Protocol) {

}

// LoginOk 登录成功的回调
func (_self *LogicProcess) LoginOk(protocol *model.Protocol) {
	log.Debugf("登录成功！:%v", protocol)
}

// LoginFail 登录失败的回调
func (_self *LogicProcess) LoginFail(protocol *model.Protocol) {

}

// ReceivedMessage 接收到消息
func (_self *LogicProcess) ReceivedMessage(protocol *model.Protocol) {
	log.Debugf("接收到服务器IM消息:%v", protocol)
	switch protocol.Type {
	case 101: //当别人申请添加自己为好友时 如果自己在线 自己就会接收到好友申请 否则就通过离线消息进行更新
		err := service.NewFriendApplyService().UpdateOne(&entity.FriendApply{
			From: util.Str2Uint64(protocol.From),
			To:   util.Str2Uint64(protocol.To),
		})
		log.Error(err)
		break
	case 102: //

	}
}

// SendOk qos中的消息发送成功 服务器成功返回
func (_self *LogicProcess) SendOk(protocol *model.Protocol) {

}
func (_self *LogicProcess) Exception(msg string) {
	log.Errorf("exception:%v", msg)
}
