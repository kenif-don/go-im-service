package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

type MessageListener interface {
	//OnReceive 当前聊天接收到消息
	OnReceive(data string)
	//OnDelMsg 如果当前聊天是 对方如果在此时删除了,会触发此通知
	OnDelMsg(data string)
	//OnSendReceive 发送的消息状态 -某消息 发送成功、发送失败
	OnSendReceive(data string)
	//OnDoChat 如果客户端停留在首页 如果有新消息进来,都会调用此接口更新最后消息和排序
	OnDoChat(data string)
	//OnFriendApply 好友申请
	OnFriendApply()
	//OnLogin 登录失效通知
	OnLogin()
	//OnLoginPwd2 输入二级密码
	OnLoginPwd2()
	//OnSend 消息发送给了服务器 但是不知道成功还是失败 只是发出去了
	OnSend(data string)
}

// SyncPutErr 同步导出
func SyncPutErr(err *utils.Error, resp *api.ResultDTOResp) []byte {
	resp.Code = uint32(api.ResultDTOCode_ERROR)
	resp.Msg = util.GetErrMsg(err)
	result, _ := proto.Marshal(resp)
	return result
}
