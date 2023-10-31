package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	"google.golang.org/protobuf/proto"
)

func InitConfig(data []byte) []byte {
	req := &api.ConfigReq{}
	resp := &api.ResultDTOResp{}
	if err := proto.Unmarshal(data, req); err != nil {
		return SyncPutErr(err, resp)
	}
	conf.InitConfig(&conf.BaseConfig{
		BaseDir:    req.BaseDir,
		LogSwitch:  req.LogSwitch.String(),
		DeviceType: req.DeviceType.String(),
		ApiHost:    req.ApiHost,
		WsHost:     req.WsHost,
	})
	resp.Code = uint32(api.ResultDTOCode_SUCCESS)
	resp.Msg = "success"
	res, _ := proto.Marshal(resp)
	return res
}
