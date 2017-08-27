package model

import (
	"time"

	"github.com/inu1255/green/config"
)

// 验证码模型
type Verify struct {
	IdItem `xorm:"extends"`
	Title  string `json:"title,omitempty" xorm:"" gev:"手机号/邮箱等"`
	Code   string `json:"code,omitempty" xorm:"not null" gev:"验证码"`
	Rest   int    `json:"rest,omitempty" xorm:"not null default 10" gev:"剩余错误次数"`
}

func (this *Verify) CanSend() bool {
	return time.Time(this.UpdateAt).Add(time.Minute).After(time.Now())
}

func (this *Verify) IsExpired() bool {
	expire_second := time.Duration(config.Cfg.Secure.CodeExpireSecond)
	return time.Time(this.UpdateAt).Add(expire_second * time.Second).Before(time.Now())
}

func (this *Verify) SendCode(title, code string) {
	config.Log.Println(title, "=>", code)
	this.Title = title
	this.Code = code
	this.Rest = config.Cfg.Secure.VerifyCodeRetry
}
