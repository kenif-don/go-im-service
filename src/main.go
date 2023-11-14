package main

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/entity"
	"fmt"
)

func main() {
	var chat *entity.Chat
	var e error
	chat, e = test()
	fmt.Println(e)
	fmt.Println(chat)
}
func test() (*entity.Chat, *utils.Error) {
	return &entity.Chat{}, nil
}
