package model

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golibs/uuid"
	"github.com/inu1255/green/config"
)

type AccessToken struct {
	Id        int    `json:"id,omitempty" xorm:"pk autoincr"`
	Token     string `json:"token,omitempty" xorm:"index" gev:"身份密钥"`
	CreateAt  Time   `json:"create_at,omitempty" xorm:"created"`
	UpdateAt  Time   `json:"-" xorm:"updated"`
	ExpiredAt Time   `json:"expired_at,omitempty" xorm:"index" gev:"过期时间"`
	UserId    int    `json:"-" xorm:""`
	Ip        string `json:"-" xorm:""`
	UA        string `json:"-" xorm:""`
	Device    string `json:"-" xorm:""`
	Uuid      string `json:"-" xorm:""`
	Action    string `json:"-" xorm:""`
}

func (this *AccessToken) ReadContextInfo(c *gin.Context) {
	if c != nil {
		UA := c.Request.Header.Get("User-Agent")
		this.Ip = c.ClientIP()
		this.UA = UA
		this.Device = c.Request.Header.Get("X-DEVICE")
		this.Uuid = c.Request.Header.Get("X-UUID")
	}
}

func NewAccessToken(user_id int, c *gin.Context) *AccessToken {
	token := &AccessToken{
		UserId:    user_id,
		Token:     uuid.Rand().Hex(),
		ExpiredAt: NewExpireTime(),
	}
	token.ReadContextInfo(c)
	return token
}

func NewExpireTime() Time {
	tokenExpire := config.Cfg.Cookie.TokenExpire
	return Time(time.Now().Add(time.Duration(tokenExpire) * time.Second))
}
