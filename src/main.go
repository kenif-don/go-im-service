package main

import (
	"IM-Service/src/util"
	"fmt"
)

func main() {
	////将wav文件读取到byte数组
	//inputFile := "D://input.wav"
	//outputFile := "D://output.wav"
	////读取文件
	//bytes, _ := os.ReadFile(inputFile)
	////变音
	//data, s, e := util.Change(bytes, 2)
	//if e != nil {
	//	return
	//}
	////时间
	//fmt.Println(s)
	//output, err := os.Create(outputFile)
	//if err != nil {
	//	panic(err)
	//}
	//defer output.Close()
	////将data写入wav文件
	//output.Write(data)
	beginIndex, endIndex := 0, 7
	//原数组
	d1 := []byte{1, 2, 3, 4, 5, 6, 7}
	//加密后的数组
	d2 := []byte{9, 9, 9}
	res := util.CoverSrcData2EnDate(d1, d2, beginIndex, endIndex)
	fmt.Println(res)
	fmt.Println(len(res))

	res = util.RevertCoveredData(res, d1[beginIndex:endIndex], beginIndex, endIndex, len(d2))
	fmt.Println(res)
	fmt.Println(len(res))
}
