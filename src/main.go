package main

import (
	"IM-Service/src/util"
	"fmt"
	"os"
)

func main() {
	//将wav文件读取到byte数组
	inputFile := "D://input.wav"
	outputFile := "D://output.wav"
	//读取文件
	bytes, _ := os.ReadFile(inputFile)
	//变音
	data, s, e := util.Change(bytes, 2)
	if e != nil {
		return
	}
	//时间
	fmt.Println(s)
	output, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	//将data写入wav文件
	output.Write(data)
}
