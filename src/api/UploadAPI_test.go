package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"

	"google.golang.org/protobuf/proto"
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
