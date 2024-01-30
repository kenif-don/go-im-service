package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

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
	//OnFriendApply 好友申请
	OnFriendApply()
	//OnLogin 登录失效通知
	OnLogin()
	//OnLoginPwd2 输入二级密码
	OnLoginPwd2()
	//OnDoVoice 收到新消息通知 用来播放语音提示
	OnDoVoice(data string)
	//OnConnectChange 连接状态变化 1-链接成功 0-链接失败
	OnConnectChange(state string)
}

// SyncPutErr 错误同步导出
func SyncPutErr(err *utils.Error, resp *api.ResultDTOResp) []byte {
	resp.Code = uint32(api.ResultDTOCode_ERROR)
	if conf.Conf.Language == conf.ZN {
		resp.Msg = err.MsgZh
	} else {
		resp.Msg = err.Msg
	}
	result, _ := proto.Marshal(resp)
	return result
}

// SyncPutSuccess 统一返回成功的封装
func SyncPutSuccess(obj interface{}, resp *api.ResultDTOResp) []byte {
	if obj != nil {
		switch obj.(type) {
		case string:
			resp.Body = obj.(string)
			break
		default:
			str, e := util.Obj2Str(obj)
			if e != nil {
				return SyncPutErr(utils.ERR_OPERATION_FAIL, resp)
			}
			resp.Body = str
		}
	}
	if resp.Body == "null" {
		resp.Body = ""
	}
	if resp.Code == 0 {
		resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	}
	if resp.Msg == "" {
		resp.Msg = "success"
	}
	result, e := proto.Marshal(resp)
	if e != nil {
		return SyncPutErr(utils.ERR_OPERATION_FAIL, resp)
	}
	return result
}
