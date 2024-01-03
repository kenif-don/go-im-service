package api

import (
	api "IM-Service/build/generated/service/v1"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/service"
	"google.golang.org/protobuf/proto"
)

// GetVersion 获取版本
func GetVersion(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	req := &api.VersionReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	result, err := service.NewVersionService().GetVersion(req.VersionCode, req.Type)
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(result, resp)
}
