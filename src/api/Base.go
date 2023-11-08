package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/util"
	"google.golang.org/protobuf/proto"
)

type MessageListener interface {
	//OnReceive 当前聊天接收到消息
	OnReceive(data []byte)
	//OnSendReceive 发送的消息状态 -某消息 发送成功、发送失败
	OnSendReceive(data []byte)
	//OnDoChats 如果客户端停留在首页 如果有新消息进来,都会调用此接口更新最后消息和排序
	OnDoChats(data []byte)
}

// SyncPutErr 同步导出
func SyncPutErr(err *utils.Error, resp *api.ResultDTOResp) []byte {
	resp.Code = uint32(api.ResultDTOCode_ERROR)
	resp.Msg = util.GetErrMsg(err)
	result, _ := proto.Marshal(resp)
	return result
}
