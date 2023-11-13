package api

import (
	api "IM-Service/build/generated/service/v1"
	"IM-Service/src/configs/log"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestUpload(t *testing.T) {
	uploadReq := &api.UploadReq{
		Path: "C:\\Users\\Administrator\\Desktop\\123.txt",
	}
	req, _ := proto.Marshal(uploadReq)
	resp := Upload(req)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
