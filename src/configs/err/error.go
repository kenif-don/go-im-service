package utils

import (
	"errors"
	"fmt"
)

var (
	ERR_USER_REGISTER_USERNAME_NULL = NewError(100, "username is null", "账号不能为空")
	ERR_USER_REGISTER_PASSWORD_NULL = NewError(101, "password is null", "密码不能为空")
	//账号长度为6-12位
	ERR_USER_REGISTER_USERNAME_LENGTH = NewError(102, "username length is error", "账号长度为6-20位")
	ERR_USER_REGISTER_PASSWORD_LENGTH = NewError(103, "password length is error", "密码长度为6-20位")
	ERR_USER_NOT_EXIST                = NewError(104, "user not exist", "用户不存在")
)

type Error struct {
	Code       int
	Msg        error
	MsgZh      error
	IsHasStack bool
}

func NewError(errCode int, msg, msgZh string) *Error {
	return &Error{
		Code:  errCode,
		Msg:   errors.New(msg),
		MsgZh: errors.New(msgZh),
	}
}

func NewSysError(err error) *Error {
	return &Error{
		Code:  500,
		Msg:   err,
		MsgZh: err,
	}
}

func (e *Error) ErrorCode() int {
	return e.Code
}

func (e *Error) Error() string {
	return fmt.Sprintf("-> Code:%d, Msg:%s, MsgZh:%s", e.Code, e.Msg.Error(), e.MsgZh.Error())
}
