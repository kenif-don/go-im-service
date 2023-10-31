package service

import (
	"IM-Service/src/configs/log"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	userService := NewUserService()
	result, err := userService.Register("test", "123456")
	if err != nil {
		log.Debug(err)
	}
	log.Debugf("%+v", result)
}
