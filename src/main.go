package main

import (
	api "IM-Service/build/generated/service/v1"
	api2 "IM-Service/src/api"
	"fmt"
)

func main() {
	resp := &api.ResultDTOResp{}
	d := api2.SyncPutSuccess("1", resp)
	fmt.Println(d)
}
