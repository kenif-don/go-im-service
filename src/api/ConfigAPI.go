package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/im"
	"IM-Service/src/service"
	"google.golang.org/protobuf/proto"
)

func InitConfig(data []byte, listener MessageListener) []byte {
	req := &api.ConfigReq{}
	resp := &api.ResultDTOResp{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	conf.InitConfig(&conf.BaseConfig{
		BaseDir:    req.BaseDir,
		LogSwitch:  req.LogSwitch.String(),
		DeviceType: req.DeviceType.String(),
		ApiHost:    req.ApiHost,
		WsHost:     req.WsHost,
	})
	service.SetListener(listener)
	im.StartIM()
	return SyncPutSuccess(nil, resp)
}
