package service

import (
	utils "IM-Service/src/configs/err"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

// GetOfflineMessage 获取离线消息
func (_self *MessageService) GetOfflineMessage() *utils.Error {
	return nil
}
