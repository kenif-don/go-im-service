package api

import (
	api "IM-Service/generated/service/v1"
	"google.golang.org/protobuf/proto"
)

type EventListener interface {
	On(data []byte)
}

// PutErr 异步导出
func PutErr(err error, resp *api.ResultDTOResp, callback EventListener) {
	resp.Code = uint32(api.ResultDTOCode_ERROR)
	resp.Msg = err.Error()
	result, _ := proto.Marshal(resp)
	callback.On(result)
}

// SyncPutErr 同步导出
func SyncPutErr(err error, resp *api.ResultDTOResp) []byte {
	resp.Code = uint32(api.ResultDTOCode_ERROR)
	resp.Msg = err.Error()
	result, _ := proto.Marshal(resp)
	return result
}
