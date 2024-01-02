package utils

import (
	"fmt"
)

var (
	ERR_PASSWORD_UPDATE_FAIL = NewError(1001, "password update fail", "密码修改失败")
	//账号长度为6-12位
	ERR_USER_USERNAME_LENGTH         = NewError(1002, "username length is error", "账号长度为6-20位")
	ERR_USER_PASSWORD_LENGTH         = NewError(1003, "password length is error", "密码长度为6-20位")
	ERR_USER_NOT_EXIST               = NewError(1004, "user not exist", "用户不存在")
	ERR_UPLOAD_FILE                  = NewError(1005, "upload file error", "上传文件失败")
	ERR_PARAM_PARSE                  = NewError(1006, "param parse error", "参数解析失败")
	ERR_GET_USER_INFO_FAIL           = NewError(1007, "get user info error", "获取用户信息失败")
	ERR_LOGIN_FAIL                   = NewError(1008, "login fail", "登录失败,请刷新失败")
	ERR_REGISTER_FAIL                = NewError(1009, "register fail", "注册失败")
	ERR_SECRET_UPDATE_FAIL           = NewError(1010, "secret update fail", "秘钥修改失败")
	ERR_NICKNAME_UPDATE_FAIL         = NewError(1011, "nickname update fail", "昵称修改失败")
	ERR_INTRO_UPDATE_FAIL            = NewError(1012, "intro update fail", "简介修改失败")
	ERR_EMAIL_UPDATE_FAIL            = NewError(1013, "email update fail", "邮箱修改失败")
	ERR_HEADIMG_UPDATE_FAIL          = NewError(1014, "headimg update fail", "头像修改失败")
	ERR_USER_UPDATE_FAIL             = NewError(1015, "user update fail", "用户信息修改失败")
	ERR_ADD_FRIEND_FAIL              = NewError(1016, "add friend fail", "添加好友失败")
	ERR_OPERATION_FAIL               = NewError(1017, "operation fail", "操作失败")
	ERR_QUERY_FAIL                   = NewError(1018, "query fail", "查询失败")
	ERR_NOT_LOGIN                    = NewError(1019, "not login", "未登录")
	ERR_SEND_FAIL_BY_NOT_TARGET      = NewError(1020, "send fail by not target", "发送失败，目标用户不存在")
	ERR_SEND_FAIL                    = NewError(1021, "send fail", "发送失败")
	ERR_OPEN_FAIL                    = NewError(1022, "open fail", "打开失败")
	ERR_ENCRYPT_FAIL                 = NewError(1023, "encrypt fail", "加密失败")
	ERR_DECRYPT_FAIL                 = NewError(1024, "decrypt fail", "解密失败")
	ERR_DEL_FAIL                     = NewError(1025, "del fail", "删除失败")
	ERR_NET_FAIL                     = NewError(1026, "net fail", "网络异常")
	ERR_NICKNAME_VALIDATE_FAIL       = NewError(1027, "nickname length is error", "昵称为1-10位")
	ERR_INTRO_VALIDATE_FAIL          = NewError(1028, "intro length is error", "简介为1-30位")
	ERR_NOT_PWD2_FAIL                = NewError(1029, "please input security password", "请输入二级密码")
	ERR_NOTIFY_FAIL                  = NewError(1030, "notify fail", "通知失败")
	ERR_RECHARGE_FAIL                = NewError(1031, "recharge fail", "充值失败")
	ERR_SELECT_PAY_NETWORK_FAIL      = NewError(1032, "please select pay network", "请选择支付网络")
	ERR_INPUT_PAY_MONEY_FAIL         = NewError(1033, "please input pay money", "请输入支付金额")
	ERR_ADD_MOOD_FAIL                = NewError(1034, "add mood fail", "添加动态失败")
	ERR_DELETE_MOOD_FAIL             = NewError(1035, "delete mood fail", "删除动态失败")
	ERR_SELECT_MOOD_FAIL             = NewError(1036, "select mood fail", "查询动态失败")
	ERR_ADD_REPLY_FAIL               = NewError(1037, "add reply fail", "回复失败")
	ERR_GET_WITHDRAWAL_FEE_FAIL      = NewError(1038, "get withdrawal fee fail", "获取提现手续费失败")
	ERR_INPUT_WITHDRAWAL_MONEY_FAIL  = NewError(1039, "please input withdrawal money", "请输入提现金额")
	ERR_INPUT_WITHDRAWAL_WALLET_FAIL = NewError(1040, "please input withdrawal wallet", "请输入提现钱包地址")
	ERR_FRIEND_GET_FAIL              = NewError(1041, "friend get fail", "获取好友失败")
	ERR_GET_USER_FAIL                = NewError(1042, "get user fail", "获取用户失败")
	ERR_MESSAGE_NOT_FOUND            = NewError(1043, "message not found", "消息不存在")
	ERR_MESSAGE_UPDATE_FAIL          = NewError(1044, "message update fail", "消息更新失败")
	ERR_GROUP_GET_FAIL               = NewError(1045, "group get fail", "获取群组失败")
	ERR_ENTER_PASSWORD               = NewError(1046, "enter password", "请输入密码")
	ERR_PASSWORD_ERROR               = NewError(1047, "password error", "密码错误")
)

type Error struct {
	Code       int
	Msg        string
	MsgZh      string
	IsHasStack bool
}

func NewError(errCode int, msg, msgZh string) *Error {
	return &Error{
		Code:  errCode,
		Msg:   msg,
		MsgZh: msgZh,
	}
}

func NewSysError(err error) *Error {
	return &Error{
		Code:  500,
		Msg:   err.Error(),
		MsgZh: err.Error(),
	}
}

func (e *Error) ErrorCode() int {
	return e.Code
}

func (e *Error) Error() string {
	return fmt.Sprintf("-> Code:%d, Msg:%s, MsgZh:%s", e.Code, e.Msg, e.MsgZh)
}
