package config

import (
	"github.com/inu1255/gev"
)

var (
	NeedLoginError        = gev.Error(1, "需要登录")
	NeedAdminError        = gev.Error(2, "需要管理员权限")
	UserNotExistError     = gev.Error(3, "用户不存在")
	UserExistError        = gev.Error(4, "账号已注册")
	PwdWrongError         = gev.Error(5, "密码错误")
	PwdLengthError        = gev.Error(6, "请输入6~32位密码")
	CodeNotSendError      = gev.Error(7, "尚未发送验证码")
	CodeExpiredError      = gev.Error(8, "验证码已过期")
	CodeNoRestError       = gev.Error(9, "验证码已失效")
	CodeWrongError        = gev.Error(10, "验证码错误")
	ImageCodeExpiredError = gev.Error(8, "图片验证码已过期")
	ImageCodeWrongError   = gev.Error(10, "图片验证码错误")
	CodeSendTooOffenError = gev.Error(11, "验证码发送太频繁")
	EmailFormatError      = gev.Error(12, "邮箱格式不正确")
	Base64FormatError     = gev.Error(13, "base64格式不正确")
	ItemNotExistError     = gev.Error(14, "不存在")
	PermissionDenied      = gev.Error(15, "没有权限")
	EmptyArrayError       = gev.Error(16, "空数组")
	NilEntityError        = gev.Error(16, "找不到对象")
	ParamLackError        = gev.Error(16, "缺少参数")
	TelphoneEmptyError    = gev.Error(16, "手机号不能为空")
	EmailEmptyError       = gev.Error(16, "邮箱不能为空")
	NeedUnbindedError     = gev.Error(16, "您已经绑定过了,不能重复绑定")
	TelphoneBindedError   = gev.Error(16, "该手机号已经绑定其它账号")
	EmailBindedError      = gev.Error(16, "该邮箱已经绑定其它账号")
	NeedTelphoneOrEmail   = gev.Error(16, "手机号或邮箱至少需要绑定一个")
)
