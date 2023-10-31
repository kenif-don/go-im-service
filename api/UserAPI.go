package api

import (
	api "IM-Service/generated/service/v1"
	"IM-Service/service"
	"google.golang.org/protobuf/proto"
)

//type RegisterListener interface {
//	On(data []byte)
//}

// func Register(data []byte, callback RegisterListener)[]byte {
func Register(data []byte) []byte {
	req := &api.RegisterReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(err, resp)
	}
	result, err := service.NewUserService().Register(req.GetUsername(), req.GetPassword())
	if err != nil {
		return SyncPutErr(err, resp)
	}
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = result.Msg
	res, _ := proto.Marshal(resp)
	return res
}
