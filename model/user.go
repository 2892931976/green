package model

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type User struct {
	OwnerItem `xorm:"extends"`
	Nickname  string `json:"nickname" xorm:"" gev:"用户昵称"`
	Telphone  string `json:"telphone" xorm:"varchar(32)" gev:"电话号码"`
	Email     string `json:"email" xorm:"varchar(128)" gev:"邮箱"`
	Password  string `json:"-" xorm:"" gev:"密码"`
	Noise     string `json:"-" xorm:"" gev:"密码加密噪音"`
	Role      string `json:"role,omitempty" xorm:"not null default 普通用户 VARCHAR(32)" gev:"用户角色"`
	Avatar    string `json:"avatar,omitempty" xorm:"" gev:"头像"`
}

func (this *User) TableName() string {
	return "user"
}

func (this *User) IsAdmin() bool {
	if this.Role == "管理员" {
		return true
	}
	return false
}

// 比较密码是否正确
func (this *User) CheckPassword(password string) bool {
	return this.Password == EncodePwd(this.Noise, password)
}

func (this *User) GetDetail(user *User) {
	if user.Avatar == "" && user.Email != "" {
		key := MD5(strings.ToLower(strings.Trim(user.Email, " \n\t")))
		user.Avatar = "https://www.gravatar.com/avatar/" + key
	}
}

// 密码加密算法
func EncodePwd(noise, password string) string {
	return MD5(noise + password)
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 登录返回数据结构
type UserAccess struct {
	*User  `xorm:"extends"`
	Access *AccessToken `json:"access,omitempty" xorm:"-"`
}

func NewUserAccess(user *User, access *AccessToken) *UserAccess {
	bean := new(UserAccess)
	bean.User = user
	bean.Access = access
	return bean
}
