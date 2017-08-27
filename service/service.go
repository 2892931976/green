package service

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev"
	"github.com/inu1255/green/config"
	"github.com/inu1255/green/model"
)

type IDataDetail interface {
	GetDetail(user *model.User)
}

type IService interface {
	Before(ctx *gin.Context) bool
	Finish(err interface{})
	After(data interface{}, err error)
}

type Service struct {
	gev.BaseService `json:"-" xorm:"-"`
	Db              *xorm.Session `json:"-" xorm:"-"`
}

// call before func
// if return false stop call func
func (this *Service) Before(ctx *gin.Context) bool {
	ok := this.BaseService.Before(ctx)
	this.Db = config.Db.NewSession()
	return ok
}

// deal with func return
// func (this *Service) After(data interface{}, err error) {}

// if panic err is the panic param
// if no panic err is nil
func (this *Service) Finish(err interface{}) {
	this.Db.Close()
	this.BaseService.Finish(err)
}

// the api will be /TagName()/xxx
func (this *Service) TagName(name string) string {
	if strings.HasSuffix(strings.ToLower(name), "service") {
		return name[:len(name)-7]
	}
	return name
}

func UserManager(param *gev.Param) {
	if param.Name == "user" && param.Type.String() == "*model.User" {
		typ := reflect.TypeOf((*model.User)(nil))
		Db := config.Db
		param.New = func(c *gin.Context) reflect.Value {
			// 当前登录用户数据
			token := c.Query("access_token")
			if c.Request.Method != http.MethodOptions && token != "" {
				now := time.Now()
				user := new(model.User)
				ok, _ := Db.Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(user)
				if ok {
					Db.Exec("update access_token set expired_at=?,update_at=? where token=?", model.NewExpireTime().String(), time.Now(), token)
					return reflect.ValueOf(user)
				}
			}
			return reflect.Zero(typ)
		}
	}
}
