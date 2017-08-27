package service

// 邮箱注册数据结构
type EmailRegisterBody struct {
	Email    string `json:"email,omitempty" xorm:"" gev:""`
	Password string `json:"password,omitempty" xorm:""`
	Code     string `json:"code,omitempty" xorm:""`
}

// 电话注册数据结构
type TelphoneRegisterBody struct {
	Telphone string `json:"telphone,omitempty" xorm:""`
	Password string `json:"password,omitempty" xorm:""`
	Code     string `json:"code,omitempty" xorm:""`
}

// 旧密码改密码
type OldpwdBody struct {
	Telphone string `json:"telphone,omitempty" xorm:""`
	Email    string `json:"email,omitempty" xorm:"" gev:""`
	Password string `json:"password,omitempty" xorm:""`
	Oldpwd   string `json:"oldpwd,omitempty" xorm:"" gev:"旧密码"`
}

type UserBody struct {
	Nickname string `json:"nickname,omitempty" xorm:"" gev:"昵称"`
	Avatar   string `json:"avatar,omitempty" xorm:"" gev:"头像"`
}
