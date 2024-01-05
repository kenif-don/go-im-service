package im

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"IM-Service/src/repository"
	"IM-Service/src/service"
	"IM-Service/src/util"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"im-sdk/client"
	"im-sdk/handler"
	"im-sdk/model"
	"strconv"
	"strings"
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
			e := conf.Conf.Client.Startup(GetLogicProcess(), websocket.New())
			if e != nil {
				if strings.Contains(e.Error(), "An existing connection was forcibly closed by the remote host") {
					log.Error(utils.ERR_NET_FAIL)
				} else {
					log.Error(e)
				}
				conf.Conf.Connected = false
			}
			if conf.Conf.Connected {
				return
			}
			time.Sleep(time.Second * 2)
		}
	}()
}

// LoginIm 长连接登录
func LoginIm() *utils.Error {
	//未登录直接返回
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return nil
	}
	loginInfo := &model.LoginInfo{
		//获取登录者 组装登录IM请求参数
		Id:     strconv.FormatUint(conf.GetLoginInfo().User.Id, 10),
		Device: conf.Base.DeviceType,
		Token:  conf.GetLoginInfo().Token,
	}
	mgr := handler.GetClientHandler().GetMessageManager()
	if mgr == nil {
		log.Error(utils.ERR_NET_FAIL)
		return log.WithError(utils.ERR_NET_FAIL)
	}
	mgr.SendLogin(loginInfo)
	return nil
}
func (_self *LogicProcess) Connected() {
	if conf.GetLoginInfo().User == nil || conf.GetLoginInfo().User.Id == 0 {
		return
	}
	//链接成功之后登录
	err := LoginIm()
	if err != nil {
		log.Error(err)
		return
	}
	if service.Listener != nil {
		service.Listener.OnConnectChange("1")
	}
}

// SendOk qos中的消息发送成功 服务器成功返回
func (_self *LogicProcess) SendOk(protocol *model.Protocol) {
	messageService := service.NewMessageService()
	messageService.UpdateReaded(protocol, 2)
	if service.Listener != nil && (protocol.Type == 1 || protocol.Type == 8) {
		//消息状态通知
		err := service.NotifySendReceive(protocol.No, 2)
		if err != nil {
			log.Error(err)
			return
		}
		var message = &entity.Message{}
		e := util.Str2Obj(protocol.Data.(string), message)
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
			chat, e = service.NewChatService().CoverChat(message.Type, util.Str2Uint64(protocol.From), false)
			if e != nil {
				log.Error(e)
				return
			}
		}
		// 通知聊天列表更新
		err = service.NewChatService().ChatNotify(chat)
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

// LoginOk 长连接登录成功的回调
func (_self *LogicProcess) LoginOk(protocol *model.Protocol) {
	conf.DiffTime = int(util.Str2Uint64(protocol.Data.(string)) - util.CurrentTime())
	service.OfflineMessageNotify()
}

// LoginFail 登录失败的回调
func (_self *LogicProcess) LoginFail(protocol *model.Protocol) {
	log.Debugf("登录失败！:%v", protocol)
}

// Logout 客户端正常退出
func (_self *LogicProcess) Logout() {
	//进行重连
	go func() {
		err := conf.Conf.Client.Reconnect(websocket.New())
		if err != nil {
			log.Error(err)
		}
	}()
}

// ReceivedMessage 接收到消息
func (_self *LogicProcess) ReceivedMessage(protocol *model.Protocol) {
	//此操作需要登录 但是当前链接未登录 直接重启
	if protocol.Ack == 500 {
		log.Error("链接报错 链接重启：%s", protocol.Data)
		err := conf.Conf.Client.Reconnect(websocket.New())
		if err == nil {
			return
		}
	}
	err := service.NewMessageService().Handler(protocol)
	if err != nil {
		log.Errorf("解析服务器IM消息失败:%v", err)
	}
}
func (_self *LogicProcess) Exception(ctx netty.ExceptionContext, e netty.Exception) {
	if service.Listener != nil {
		service.Listener.OnConnectChange("0")
	}
	log.Error(e)
	conf.Conf.Connected = false
	log.Debug("服务器断开连接,进行重连")
	go func() {
		for {
			err := conf.Conf.Client.Reconnect(websocket.New())
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(5 * time.Second)
		}
	}()
}

// HandleEvent 处理事件
func (_self *LogicProcess) HandleEvent(ctx netty.EventContext, event netty.Event) {
	log.Debug(event)
}
