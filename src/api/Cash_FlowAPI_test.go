package api

import (
	api "go-im-service/build/generated/service/v1"
	"go-im-service/src/configs/log"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestPagingCashFlow(t *testing.T) {
	resp := PagingCashFlow(1)
	result := &api.ResultDTOResp{}
	err := proto.Unmarshal(resp, result)
	if err != nil {
		log.Error(err)
	}
	log.Debug(result)
}
