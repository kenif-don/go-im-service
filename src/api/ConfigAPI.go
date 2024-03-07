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

func SelectConfig() []byte {
	resp := &api.ResultDTOResp{}
	res, err := service.SelectConfig()
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(res, resp)
}

// SetLanguage 设置语言
func SetLanguage(language string) []byte {
	service.SetLanguage(language)
	resp := &api.ResultDTOResp{}
	return SyncPutSuccess(nil, resp)
}
func GetAgent() []byte {
	resp := &api.ResultDTOResp{}
	res, err := service.GetAgent()
	if err != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	return SyncPutSuccess(res, resp)
}
