package main

import (
	"fmt"
	"go-im-service/src/util"
)

func main() {
	data, err := util.DecryptAes("OIa6SS1TNS0GWoMSHPcaWPd7eh8w3WXjXwCkTmPPIHU=", util.MD5("safe_111026"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
