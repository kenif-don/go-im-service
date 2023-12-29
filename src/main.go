package main

import (
	"IM-Service/src/util"
	"fmt"
)

func main() {
	endwith, e := util.GetFileType("C:\\Users\\Administrator\\Desktop\\cbd3480f306239cf12a957b656d185e5.txt", make([]byte, 0))
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(endwith)
}
