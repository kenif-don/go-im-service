package main

import (
	"IM-Service/src/util"
	"fmt"
)

func main() {
	endWith, e := util.GetFileType("C:\\Users\\Administrator\\Desktop\\bug.txt", make([]byte, 0))
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(endWith)
}
