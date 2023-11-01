package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"google.golang.org/protobuf/proto"
)

type EventListener interface {
	On(data []byte)
}

// PutErr 异步导出
func PutErr(err *utils.Error, resp *api.ResultDTOResp, callback EventListener) {
	resp.Code = uint32(api.ResultDTOCode_ERROR)
	resp.Msg = err.Msg
	result, _ := proto.Marshal(resp)
	callback.On(result)
}

// SyncPutErr 同步导出
func SyncPutErr(err *utils.Error, resp *api.ResultDTOResp) []byte {
	resp.Code = uint32(api.ResultDTOCode_ERROR)
	resp.Msg = err.MsgZh
	result, _ := proto.Marshal(resp)
	return result
}
