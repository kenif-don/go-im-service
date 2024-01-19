package main

import (
	"IM-Service/src/util"
	"fmt"
)

func main() {
	data, err := util.DecryptAes("OIa6SS1TNS0GWoMSHPcaWPd7eh8w3WXjXwCkTmPPIHU=", util.MD5("safe_111026"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
