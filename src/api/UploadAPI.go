package api

import (
	api "go-im-service/build/generated/service/v1"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/service"
	"go-im-service/src/util"

	"google.golang.org/protobuf/proto"
)

func Upload(data []byte) []byte {
	resp := &api.ResultDTOResp{}
	if !service.ValidatePwd2() {
		return SyncPutErr(utils.ERR_NOT_PWD2_FAIL, resp)
	}
	req := &api.UploadReq{}
	if e := proto.Unmarshal(data, req); e != nil {
		return SyncPutErr(utils.ERR_PARAM_PARSE, resp)
	}
	url, err := util.Upload(req.Path, "")
	if err != nil {
		return SyncPutErr(err, resp)
	}
	return SyncPutSuccess(url, resp)
}
